package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hpcloud/tail"
	"github.com/humio/cli/api"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

type eventList struct {
	Type     string            `json:"type"`
	Fields   map[string]string `json:"fields"`
	Messages []string          `json:"messages"`
}

type lineHandler interface {
	handleLine(line string)
}

type multiLineHandlerMode int

const (
	multiLineHandlerModeBeginsWith multiLineHandlerMode = iota
	multiLineHandlerModeContinuesWith
)

type multiLineHandler struct {
	lineHandler lineHandler
	regex       *regexp.Regexp
	mode        multiLineHandlerMode
	buf         bytes.Buffer
}

func (h *multiLineHandler) handleLine(line string) {
	isMatch := h.regex.MatchString(line)

	switch h.mode {
	case multiLineHandlerModeBeginsWith:
		if isMatch {
			fullLine := h.buf.String()
			h.buf.Reset()
			h.lineHandler.handleLine(fullLine)
		}

	case multiLineHandlerModeContinuesWith:
		if !isMatch {
			fullLine := h.buf.String()
			h.buf.Reset()
			h.lineHandler.handleLine(fullLine)
		}
	}
	h.buf.WriteString(line)
	h.buf.WriteString("\n")
}

func tailFile(filepath string, quiet bool, seekToEnd bool, handler lineHandler) error {
	tailConfig := tail.Config{Follow: true}
	if seekToEnd {
		tailConfig.Location = &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}
	}

	t, err := tail.TailFile(filepath, tailConfig)

	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		handler.handleLine(line.Text)
		if !quiet {
			fmt.Println(line.Text)
		}
	}

	waitForInterrupt()

	err = t.Stop()

	return err
}

func streamStdin(repo string, quiet bool, handler lineHandler) error {
	log.Println("Humio Attached to StdIn, Forwarding to '" + repo + "'")

	var reader io.Reader = os.Stdin
	if !quiet {
		reader = io.TeeReader(reader, os.Stdout)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		handler.handleLine(scanner.Text())
	}

	return scanner.Err()
}

func waitForInterrupt() {
	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// This goroutine executes a blocking receive for
	// signals. When it gets one it'll print it out
	// and then notify the program that it can finish.
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		close(done)
	}()

	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	<-done
}

type logSenderErrorBehaviour int

const (
	logSenderErrorBehaviourDrop logSenderErrorBehaviour = iota
	logSenderErrorBehaviourCrash
)

type logSender struct {
	apiClient           *api.Client
	url                 string
	fields              map[string]string
	parserName          string
	events              chan string
	finishedSending     chan struct{}
	maxAttemptsPerBatch int
	errorBehaviour      logSenderErrorBehaviour
	batchSizeLines      int
	batchSizeBytes      int
	batchTimeout        time.Duration
}

func (s *logSender) handleLine(line string) {
	s.events <- line
}

func (s *logSender) finish() {
	close(s.events)
	<-s.finishedSending
}

func (s *logSender) start() {
	go func() {
		defer func() { close(s.finishedSending) }()
		batch := make([]string, 0, s.batchSizeLines)

		for {
			bytes := 0

			e, more := <-s.events
			if !more {
				break
			}

			batch = append(batch, e)
			bytes += len(e)

			timeout := time.After(s.batchTimeout)

			if s.batchSizeBytes > 0 && bytes > s.batchSizeBytes {
				goto send
			}

		loop:
			for {
				select {
				case e, more := <-s.events:
					if !more {
						break loop
					}
					batch = append(batch, e)
					bytes += len(e)
					if len(batch) >= s.batchSizeLines || (s.batchSizeBytes > 0 && bytes > s.batchSizeBytes) {
						break loop
					}
				case <-timeout:
					break loop
				}
			}

		send:
			s.sendBatch(batch)

			batch = batch[:0]
			bytes = 0
		}
		if len(batch) > 0 {
			s.sendBatch(batch)
		}
	}()
}

func (s *logSender) sendBatch(messages []string) {
	ship := func() error {
		var eg errgroup.Group

		pr, pw := io.Pipe()

		jsonBody := []eventList{{
			Type:     s.parserName,
			Fields:   s.fields,
			Messages: messages,
		}}

		eg.Go(func() error {
			defer pw.Close()
			return json.NewEncoder(pw).Encode(jsonBody)
		})

		var resp *http.Response

		eg.Go(func() error {
			var err error
			resp, err = s.apiClient.HTTPRequest(http.MethodPost, s.url, pr)
			return err
		})

		err := eg.Wait()

		if err != nil {
			return err
		}

		if resp.StatusCode > 400 {
			responseData, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				return fmt.Errorf("error reading http response body: %w", err)
			}

			return fmt.Errorf("bad response while sending events (status='%s'): %s", resp.Status, responseData)
		} else {
			// discard the response in order to re-use the connection
			_, _ = io.Copy(ioutil.Discard, resp.Body)
			_ = resp.Body.Close()
		}

		return nil
	}

	var err error
	for i := 0; i < s.maxAttemptsPerBatch; i++ {
		if i > 0 {
			backOff := time.Duration(0.5*math.Pow(2, float64(i-1))*1000) * time.Millisecond
			log.Printf("Backoff for %v...", backOff)
			time.Sleep(backOff)
		}
		err = ship()
		if err == nil {
			break
		}
		log.Printf("Error while sending logs to Humio. Retrying %d more times. Error message: %v", s.maxAttemptsPerBatch-i-1, err)
	}

	if err != nil {
		switch s.errorBehaviour {
		case logSenderErrorBehaviourCrash:
			log.Fatalf("Error sending logs to Humio: %v", err)
		case logSenderErrorBehaviourDrop:
			log.Printf("Error sending logs to Humio, dropping %d events: %v", len(messages), err)
		}
	}
}

func newIngestCmd() *cobra.Command {
	var parserName, filepath, label, ingestToken, multiLineBeginsWith, multiLineContinuesWith, fieldsJson string
	var openBrowser, noSession, quiet, failOnError, tailSeekToEnd bool
	var retries, batchSizeLines, batchSizeBytes, batchTimeoutMs int

	cmd := cobra.Command{
		Use:   "ingest [flags] [<repo>]",
		Short: "Send data to Humio.",
		Long: `Listens to stdin and sends all input to the repository <repo>.
If the --ingest-token flag is specified, the repo associated with the ingest token will be used.
Otherwise, if <repo> is not specified, Humio will use your 'sandbox' repository as
destination.

It can be handy to specify the parser to be used to ingest the
data on arrival - i.e. the type of data you are sending.
The value of --parser will take precedence over any parser that
is associated with the ingest token set by --ingest-token.

You can pipe the output of another process through humio:

  $ tail -f /var/log/syslog | humio ingest --ingest-token=af21... --parser=syslog

Alternatively, you can use the --tail=<file> argument, which
has the same effect.`,
		ValidArgs: []string{"repo"},
		Args:      cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var repo string

			// Default to sending to the sandbox
			if l := len(args); l == 1 {
				repo = args[0]
			} else {
				repo = "sandbox"
			}

			var opts []func(config *api.Config)
			if ingestToken != "" {
				opts = append(opts, func(config *api.Config) {
					config.Token = ingestToken
				})
			}

			client := NewApiClient(cmd, opts...)

			var key string
			fields := map[string]string{}

			if fieldsJson != "" {
				var f map[string]string
				err := json.Unmarshal([]byte(fieldsJson), &f)
				if err != nil {
					log.Fatalf("Error parsing --fields-json value: %v", err)
				}
				for k := range f {
					fields[k] = f[k]
				}
			} else {
				if !noSession {
					u, _ := uuid.NewV4()
					sessionID := u.String()
					fields["@session"] = sessionID

					if label == "" && !noSession {
						key = "%40session%3D" + sessionID
					}
				}
				if label != "" {
					fields["@label"] = label
					key = "%40label%3D" + label
				}
			}

			// Open the browser (First so it has a chance to load)
			if openBrowser {
				browserURL, err := client.Address().Parse(fmt.Sprintf("/%s/search?live=true&start=1d&query=%s", repo, key))
				if err != nil {
					fmt.Println(fmt.Errorf("could not open browser: %v", err))
				}
				err = open.Start(browserURL.String())
				if err != nil {
					fmt.Println(fmt.Errorf("could not open browser: %v", err))
				}
			}

			var url string
			if ingestToken != "" {
				url = "/api/v1/ingest/humio-unstructured"
			} else {
				url = "/api/v1/repositories/" + repo + "/ingest-messages"
			}

			sender := logSender{
				apiClient:           client,
				url:                 url,
				fields:              fields,
				parserName:          parserName,
				maxAttemptsPerBatch: retries + 1,
				events:              make(chan string, batchSizeLines),
				finishedSending:     make(chan struct{}),
				batchSizeLines:      batchSizeLines,
				batchSizeBytes:      batchSizeBytes,
				batchTimeout:        time.Duration(batchTimeoutMs) * time.Millisecond,
			}

			if failOnError {
				sender.errorBehaviour = logSenderErrorBehaviourCrash
			}

			sender.start()

			var lineHandler lineHandler = &sender

			switch {
			case multiLineBeginsWith != "" && multiLineContinuesWith != "":
				log.Fatalf("Cannot specify both --multiline-begins-with and --multiline-continues-with")
			case multiLineBeginsWith != "":
				lineHandler = &multiLineHandler{
					lineHandler: lineHandler,
					regex:       regexp.MustCompile(multiLineBeginsWith),
					mode:        multiLineHandlerModeBeginsWith,
				}
			case multiLineContinuesWith != "":
				lineHandler = &multiLineHandler{
					lineHandler: lineHandler,
					regex:       regexp.MustCompile(multiLineContinuesWith),
					mode:        multiLineHandlerModeContinuesWith,
				}
			}

			var err error
			if filepath != "" {
				err = tailFile(filepath, quiet, tailSeekToEnd, lineHandler)
			} else {
				err = streamStdin(repo, quiet, lineHandler)
			}

			sender.finish()

			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&parserName, "parser", "p", "default", "Use a specific parser for ingestion.")
	cmd.Flags().StringVarP(&filepath, "tail", "f", "", "A file to tail instead of listening to stdin.")
	cmd.Flags().BoolVarP(&tailSeekToEnd, "tail-end", "E", false, "When used with --tail, start from the end of the file and follow it. Equivalent to 'tail -f -n0 <file>'")
	cmd.Flags().StringVarP(&ingestToken, "ingest-token", "i", "", "Use the specified ingest token instead of the API token.")
	cmd.Flags().BoolVarP(&openBrowser, "open", "o", false, "Open the browser with live tail of the stream.")
	cmd.Flags().StringVarP(&label, "label", "l", "", "Adds a @label=<lavel> field to each event. This can help you find specific data send by the CLI when searching in the UI.")
	cmd.Flags().BoolVarP(&noSession, "no-session", "n", false, "No @session field will be added to each event. @session assigns a new UUID to each executing of the Humio CLI.")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Don't print ingested data to stdout.")
	cmd.Flags().BoolVarP(&failOnError, "fail", "e", false, "Stop processing more input when sending events has failed (after the allowed number of retries).")
	cmd.Flags().IntVarP(&retries, "retries", "r", 2, "Number of retries when Humio sending events.")
	cmd.Flags().IntVarP(&batchSizeLines, "batch-lines", "L", 500, "Max number of events to send in one batch.")
	cmd.Flags().IntVarP(&batchSizeBytes, "batch-bytes", "B", 1024*1024, "Max number of bytes to send in one batch.")
	cmd.Flags().IntVarP(&batchTimeoutMs, "batch-timeout", "T", 100, "Max duration in milliseconds to wait before sending an incomplete batch.")
	cmd.Flags().StringVarP(&multiLineBeginsWith, "multiline-begins-with", "", "", "Operate in multi line mode. Each multi line event starts with the specified regexp pattern.")
	cmd.Flags().StringVarP(&multiLineContinuesWith, "multiline-continues-with", "", "", "Operate in multi line mode. Each multi line event is continued with the specified regexp pattern.")
	cmd.Flags().StringVarP(&fieldsJson, "fields-json", "J", "", "Add the supplied json object to each object as structured fields.")

	return &cmd
}
