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

	"github.com/hpcloud/tail"
	uuid "github.com/satori/go.uuid"
	cli "gopkg.in/urfave/cli.v2"
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

func tailFile(server server, filepath string) {

	// Join Tail

	t, err := tail.TailFile(filepath, tail.Config{Follow: true})

	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		sendLine(server, line.Text)
	}

	tailError := t.Wait()

	if tailError != nil {
		log.Fatal(tailError)
	}
}

func streamStdin(server server) {
	log.Println("Humio Attached to StdIn")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		sendLine(server, text)
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

func Ingest(c *cli.Context) error {
	filepath := c.Args().Get(0)

	name := ""
	if c.String("name") == "" && filepath != "" {
		name = filepath
	} else {
		name = c.String("name")
	}

	var parserName string
	if c.String("parser") != "" {
		parserName = c.String("parser")
	} else {
		parserName = "default"
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

	fields := map[string]string{
		"@session": sessionID,
	}

	if name != "" {
		fields["@name"] = name
	}

	startSending(config, fields, parserName)

	if filepath != "" {
		tailFile(config, filepath)
	} else {
		streamStdin(config)
	}
	return nil
}

func startSending(server server, fields map[string]string, parserName string) {
	go func() {
		var batch []string
		for {
			select {
			case v := <-events:
				batch = append(batch, v)
				if len(batch) >= batchLimit {
					sendBatch(server, batch, fields, parserName)
					batch = batch[:0]
				}
			default:
				if len(batch) > 0 {
					sendBatch(server, batch, fields, parserName)
					batch = batch[:0]
				}
				// Avoid busy waiting
				batch = append(batch, <-events)
			}
		}
	}()
}

func sendLine(server server, line string) {
	events <- line
}

func sendBatch(server server, messages []string, fields map[string]string, parserName string) {

	lineJSON, marshalErr := json.Marshal([1]eventList{
		eventList{
			Type:     parserName,
			Fields:   fields,
			Messages: messages,
		}})

	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	ingestURL := server.URL + "api/v1/repositories/" + server.Repo + "/ingest-messages"
	lineReq, reqErr := http.NewRequest("POST", ingestURL, bytes.NewBuffer(lineJSON))
	lineReq.Header.Set("Authorization", "Bearer "+server.Token)
	lineReq.Header.Set("Content-Type", "application/json")
	check(reqErr)

	resp, clientErr := client.Do(lineReq)
	check(clientErr)

	if resp.StatusCode > 400 {
		responseData, readErr := ioutil.ReadAll(resp.Body)
		check(readErr)
		log.Fatal(string(responseData))
	}

	resp.Body.Close()
}
