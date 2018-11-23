package command

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hpcloud/tail"
	"github.com/humio/cli/api"
	uuid "github.com/satori/go.uuid"
	"github.com/skratchdot/open-golang/open"
)

var batchLimit = 500
var events = make(chan string, batchLimit)

type eventList struct {
	Type     string            `json:"type"`
	Fields   map[string]string `json:"fields"`
	Messages []string          `json:"messages"`
}

type event struct {
	RawString string `json:"rawstring"`
}

func (f *IngestCommand) tailFile(client *api.Client, repo string, filepath string) {

	// Join Tail

	t, err := tail.TailFile(filepath, tail.Config{Follow: true})

	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		f.sendLine(line.Text)
		fmt.Println(line.Text)
	}

	tailError := t.Wait()

	if tailError != nil {
		log.Fatal(tailError)
	}
}

func (f *IngestCommand) streamStdin(repo string) {
	log.Println("Humio Attached to StdIn, Forwarding to '" + repo + "'")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		f.sendLine(text)
		// TODO: We should be able to do this more efficiently.
		// Somehow connecting Stdin to Stdout
		fmt.Println(text)
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

type IngestCommand struct {
	Meta
}

func (f *IngestCommand) Help() string {
	helpText := `
Usage: humio ingest [-token=<ingest-token>] [-parser=<parser>] [<repo>]

	Listens to stdout and sends all output to the repository <repo>.
	If <repo> is not specified, humio will use your 'sandbox' repository as
	destination.

	It can be handy to specify the parser to be used to ingest the
	data on arrival - i.e. the type of data you are sending.

	The value of -parser will take precedence over any parser that
	is associated with the ingest token set by -token.

  You can pipe the output of another process through humio:

    $ tail -f /var/log/syslog | humio ingest -token=af21... -parser=syslog

	Alternatively, you can use the -tail=<file> argument, which
	has the same effect.

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *IngestCommand) Synopsis() string {
	return "Send data to Humio."
}

func (f *IngestCommand) Name() string { return "ingest" }

func (f *IngestCommand) Run(args []string) int {
	var parserName, filepath, label string
	var openBrowser, noSession bool

	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }
	flags.StringVar(&parserName, "parser", "kv", "Use a specific parser for ingestion.")
	flags.StringVar(&filepath, "tail", "", "A file to tail instead of listening to stdout.")
	flags.BoolVar(&openBrowser, "open", false, "Open the browser with live tail of the stream.")
	flags.StringVar(&label, "label", "", "Add @label=<value> on all events. Making it easy to search for.")
	flags.BoolVar(&noSession, "no-session", false, "Don't add a @session field to events.")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check that we got one argument
	args = flags.Args()
	if l := len(args); l >= 2 {
		f.Ui.Error("This command takes zero or one argument: <repo>")
		f.Ui.Error(commandErrorText(f))
		return 1
	}

	var repo string
	if l := len(args); l == 1 {
		repo = args[0]
	} else {
		repo = "sandbox"
	}

	client, err := f.Meta.Client()
	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

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
		open.Start(client.Address() + repo + "/search?live=true&start=1d&query=" + key)
	}

	f.startSending(client, repo, fields, parserName)

	if filepath != "" {
		f.tailFile(client, repo, filepath)
	} else {
		f.streamStdin(repo)
	}

	return 0
}

func (f *IngestCommand) startSending(client *api.Client, repo string, fields map[string]string, parserName string) {
	go func() {
		var batch []string
		for {
			select {
			case v := <-events:
				batch = append(batch, v)
				if len(batch) >= batchLimit {
					f.sendBatch(client, repo, batch, fields, parserName)
					batch = batch[:0]
				}
			default:
				if len(batch) > 0 {
					f.sendBatch(client, repo, batch, fields, parserName)
					batch = batch[:0]
				}
				// Avoid busy waiting
				batch = append(batch, <-events)
			}
		}
	}()
}

func (f *IngestCommand) sendLine(line string) {
	events <- line
}

func (f *IngestCommand) sendBatch(client *api.Client, repo string, messages []string, fields map[string]string, parserName string) {
	lineJSON, err := json.Marshal([1]eventList{
		eventList{
			Type:     parserName,
			Fields:   fields,
			Messages: messages,
		}})

	if err != nil {
		fmt.Printf("error while sending data: %v", err)
		return
	}

	url := "api/v1/repositories/" + repo + "/ingest-messages"
	resp, err := client.HttpPOST(url, bytes.NewBuffer(lineJSON))

	if err != nil {
		f.Ui.Error(fmt.Sprintf("error while sending data: %v", err))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		responseData, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			f.Ui.Error(fmt.Sprintf("error while sending data: %v", err))
			return
		}

		f.Ui.Error(fmt.Sprintf("Bad response while sending events: %s", string(responseData)))
	}
}
