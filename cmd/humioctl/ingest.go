package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/humio/cli/shipper"

	"github.com/gofrs/uuid"
	"github.com/hpcloud/tail"
	"github.com/humio/cli/api"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func tailFile(cmd *cobra.Command, filepath string, quiet bool, seekToEnd bool, handler shipper.LineHandler) error {
	tailConfig := tail.Config{Follow: true}
	if seekToEnd {
		tailConfig.Location = &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}
	}

	t, err := tail.TailFile(filepath, tailConfig)

	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		handler.HandleLine(line.Text)
		if !quiet {
			fmt.Fprintln(cmd.OutOrStdout(), line.Text)
		}
	}

	waitForInterrupt()

	err = t.Stop()

	return err
}

func streamStdin(repo string, quiet bool, handler shipper.LineHandler) error {
	log.Println("Humio Attached to StdIn, Forwarding to '" + repo + "'")

	var reader io.Reader = os.Stdin
	if !quiet {
		reader = io.TeeReader(reader, os.Stdout)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		handler.HandleLine(scanner.Text())
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
		log.Println()
		log.Println(sig)
		close(done)
	}()

	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	<-done
}

func newIngestCmd() *cobra.Command {
	var parserName, filepath, label, ingestToken, multiLineBeginsWith, multiLineContinuesWith, fieldsJson string
	var openBrowser, noSession, quiet, failOnError, tailSeekToEnd bool
	var retries, batchSizeLines, batchSizeBytes, batchTimeoutMs int

	cmd := cobra.Command{
		Use:   "ingest [flags] repo",
		Short: "Send data to Humio.",
		Long: `Listens to stdin and sends all input to the repository <repo>.
If the --ingest-token flag is specified, the repo associated with the ingest token will be used.

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

			if l := len(args); l == 1 {
				repo = args[0]
			} else {
				log.Fatal("Must specify repo to ingest data")
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

			// Open the browser
			if openBrowser {
				browserURL, err := client.Address().Parse(fmt.Sprintf("%s/search?live=true&start=1d&query=%s", repo, key))
				if err != nil {
					cmd.PrintErrf("Could not parse url: %v\n", err)
				}
				err = open.Start(browserURL.String())
				if err != nil {
					cmd.PrintErrf("Could not open browser: %v\n", err)
				}
			}

			var url string
			if ingestToken != "" {
				url = "api/v1/ingest/humio-unstructured"
			} else {
				url = "api/v1/repositories/" + repo + "/ingest-messages"
			}

			sender := shipper.LogShipper{
				APIClient:           client,
				URL:                 url,
				Fields:              fields,
				ParserName:          parserName,
				MaxAttemptsPerBatch: retries + 1,
				BatchSizeLines:      batchSizeLines,
				BatchSizeBytes:      batchSizeBytes,
				BatchTimeout:        time.Duration(batchTimeoutMs) * time.Millisecond,
				Logger:              log.New(cmd.ErrOrStderr(), "", log.LstdFlags).Printf,
			}

			if failOnError {
				sender.ErrorBehaviour = shipper.ErrorBehaviourPanic
			}

			sender.Start()

			var lineHandler shipper.LineHandler = &sender

			switch {
			case multiLineBeginsWith != "" && multiLineContinuesWith != "":
				log.Fatalf("Cannot specify both --multiline-begins-with and --multiline-continues-with")
			case multiLineBeginsWith != "":
				lineHandler = &shipper.MultiLineHandler{
					LineHandler: lineHandler,
					Regex:       regexp.MustCompile(multiLineBeginsWith),
					Mode:        shipper.MultiLineHandlerModeBeginsWith,
				}
			case multiLineContinuesWith != "":
				lineHandler = &shipper.MultiLineHandler{
					LineHandler: lineHandler,
					Regex:       regexp.MustCompile(multiLineContinuesWith),
					Mode:        shipper.MultiLineHandlerModeContinuesWith,
				}
			}

			var err error
			if filepath != "" {
				err = tailFile(cmd, filepath, quiet, tailSeekToEnd, lineHandler)
			} else {
				err = streamStdin(repo, quiet, lineHandler)
			}

			sender.Finish()

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
	cmd.Flags().StringVarP(&label, "label", "l", "", "Adds a @label=<label> field to each event. This can help you find specific data sent by the CLI when searching in the UI.")
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
