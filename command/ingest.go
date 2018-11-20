package command

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
	"time"

	"github.com/hpcloud/tail"
	uuid "github.com/satori/go.uuid"
	cli "gopkg.in/urfave/cli.v2"
)

var client = &http.Client{}

type server struct {
	URL   string
	Token string
	Repo  string
}

func getServerConfig(c *cli.Context) (server, error) {
	config := server{
		Repo:  c.String("repo"),
		Token: c.String("token"),
		URL:   c.String("url"),
	}
	return config, nil
}

var batchLimit = 500
var events = make(chan event, 500)

type eventList struct {
	Tags   map[string]string `json:"tags"`
	Events []event           `json:"events"`
}

type event struct {
	Timestamp  string            `json:"timestamp"`
	Attributes map[string]string `json:"attributes"`
	RawString  string            `json:"rawstring"`
}

func tailFile(server server, name string, sessionID string, filepath string) {

	// Join Tail

	t, err := tail.TailFile(filepath, tail.Config{Follow: true})

	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		sendLine(server, name, sessionID, line.Text)
	}

	tailError := t.Wait()

	if tailError != nil {
		log.Fatal(tailError)
	}
}

func streamStdin(server server, name string, sessionID string) {
	log.Println("Humio Attached to StdIn")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sendLine(server, name, sessionID, scanner.Text())
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

func Ingest(c *cli.Context) error {
	filepath := c.Args().Get(0)

	name := ""
	if c.String("name") == "" && filepath != "" {
		name = filepath
	} else {
		name = c.String("name")
	}

	u, _ := uuid.NewV4()

	sessionID := u.String()

	config, _ := getServerConfig(c)

	ensureToken(config)

	// Open the browser (First so it has a chance to load)
	// key := ""
	// if name == "" {
	// 	key = "%40session%3D" + server.SessionID
	// } else {
	// 	key = "%40name%3D" + server.Name
	// }
	// open.Run(server.ServerURL + server.RepoID + "/search?live=true&start=1d&query=" + key)

	startSending(config)

	if filepath != "" {
		tailFile(config, name, sessionID, filepath)
	} else {
		streamStdin(config, name, sessionID)
	}
	return nil
}

func startSending(server server) {
	go func() {
		var batch []event
		for {
			select {
			case v := <-events:
				batch = append(batch, v)
				if len(batch) >= batchLimit {
					sendBatch(server, batch)
					batch = batch[:0]
				}
			default:
				if len(batch) > 0 {
					sendBatch(server, batch)
					batch = batch[:0]
				}
				// Avoid busy waiting
				batch = append(batch, <-events)
			}
		}
	}()
}

func sendLine(server server, name string, sessionID string, line string) {
	theEvent := event{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Attributes: map[string]string{
			"@session": sessionID,
			"@name":    name,
		},
		RawString: line,
	}

	events <- theEvent
}

func sendBatch(server server, events []event) {

	lineJSON, marshalErr := json.Marshal([1]eventList{
		eventList{
			Tags:   map[string]string{},
			Events: events,
		}})

	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	ingestURL := server.URL + "/api/v1/dataspaces/" + server.Repo + "/ingest"
	lineReq, reqErr := http.NewRequest("POST", ingestURL, bytes.NewBuffer(lineJSON))
	lineReq.Header.Set("Authorization", "Bearer "+server.Token)
	lineReq.Header.Set("Content-Type", "application/json")

	if reqErr != nil {
		panic(reqErr)
	}

	resp, clientErr := client.Do(lineReq)
	if clientErr != nil {
		panic(clientErr)
	}
	if resp.StatusCode > 400 {
		responseData, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
		log.Fatal(string(responseData))
	}
	resp.Body.Close()
}
