package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofrs/uuid"
	"github.com/hpcloud/tail"
	"github.com/humio/cli/api"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var batchLimit = 500
var events = make(chan string, batchLimit)

type eventList struct {
	Type     string            `json:"type"`
	Fields   map[string]string `json:"fields"`
	Messages []string          `json:"messages"`
}

func tailFile(client *api.Client, repo string, filepath string, quiet bool) {

	// Join Tail

	t, err := tail.TailFile(filepath, tail.Config{Follow: true})

	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		sendLine(line.Text)
		if !quiet {
			fmt.Println(line.Text)
		}
	}

	tailError := t.Wait()

	if tailError != nil {
		log.Fatal(tailError)
	}
}

func streamStdin(repo string, quiet bool) {
	log.Println("Humio Attached to StdIn, Forwarding to '" + repo + "'")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		sendLine(text)
		// TODO: We should be able to do this more efficiently.
		// Somehow connecting Stdin to Stdout
		if !quiet {
			fmt.Println(text)
		}
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	waitForInterrupt()
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
		done <- true
	}()

	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	<-done
}

func startSending(client *api.Client, repo string, fields map[string]string, parserName string) {
	go func() {
		var batch []string
		for {
			select {
			case v := <-events:
				batch = append(batch, v)
				if len(batch) >= batchLimit {
					sendBatch(client, repo, batch, fields, parserName)
					batch = batch[:0]
				}
			default:
				if len(batch) > 0 {
					sendBatch(client, repo, batch, fields, parserName)
					batch = batch[:0]
				}
				// Avoid busy waiting
				batch = append(batch, <-events)
			}
		}
	}()
}

func sendLine(line string) {
	events <- line
}

func sendBatch(client *api.Client, repo string, messages []string, fields map[string]string, parserName string) {
	lineJSON, err := json.Marshal([1]eventList{
		{
			Type:     parserName,
			Fields:   fields,
			Messages: messages,
		}})

	if err != nil {
		fmt.Printf("error while sending data: %v", err)
		return
	}

	url := "api/v1/repositories/" + repo + "/ingest-messages"
	resp, err := client.HTTPRequest(http.MethodPost, url, bytes.NewBuffer(lineJSON))

	if err != nil {
		fmt.Println((fmt.Errorf("error while sending data: %v", err)))
	}

	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		responseData, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(fmt.Errorf("error while sending data: %v", err))
		}

		fmt.Println((fmt.Errorf("Bad response while sending events: %s", string(responseData))))
	}
}

func newIngestCmd() *cobra.Command {
	var parserName, filepath, label string
	var openBrowser, noSession, quiet bool

	cmd := cobra.Command{
		Use:   "ingest [flags] [<repo>]",
		Short: "Send data to Humio.",
		Long: `Listens to stdin and sends all input to the repository <repo>.
If <repo> is not specified, Humio will use your 'sandbox' repository as
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

			client := NewApiClient(cmd)

			var key string
			fields := map[string]string{}

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

			// Open the browser (First so it has a chance to load)
			if openBrowser {
				err := open.Start(client.Address() + repo + "/search?live=true&start=1d&query=" + key)
				if err != nil {
					fmt.Println(fmt.Errorf("could not open browser: %v", err))
				}
			}

			startSending(client, repo, fields, parserName)

			if filepath != "" {
				tailFile(client, repo, filepath, quiet)
			} else {
				streamStdin(repo, quiet)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&parserName, "parser", "p", "default", "Use a specific parser for ingestion.")
	cmd.Flags().StringVarP(&filepath, "tail", "f", "", "A file to tail instead of listening to stdin.")
	cmd.Flags().StringP("ingest-token", "i", "", "The ingest token to use. Defaults to your Account API token.")
	cmd.Flags().BoolVarP(&openBrowser, "open", "o", false, "Open the browser with live tail of the stream.")
	cmd.Flags().StringVarP(&label, "label", "l", "", "Adds a @label=<lavel> field to each event. This can help you find specific data send by the CLI when searching in the UI.")
	cmd.Flags().BoolVarP(&noSession, "no-session", "n", false, "No @session field will be added to each event. @session assigns a new UUID to each executing of the Humio CLI.")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Don't print ingested data to stdout.")

	return &cmd
}
