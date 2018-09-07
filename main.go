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
	"time"

	"github.com/hpcloud/tail"
	"github.com/satori/go.uuid"
	// "github.com/skratchdot/open-golang/open"
	"gopkg.in/urfave/cli.v2"
)


////////////////////////////////////////////////////////////////////////////////
///// Globals //////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// This is set by the release script.
var version = "master"
var client = &http.Client{}


////////////////////////////////////////////////////////////////////////////////
///// main function ////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func main() {
	app := &cli.App{
		Name:  "humio",
		Usage: "humio [options] <filepath>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "Your Humio API Token",
				EnvVars: []string{"HUMIO_API_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "dataspace",
				Aliases: []string{"d"},
				Value:   "scratch",
				Usage:   "The dataspace to stream to. Defaults to your scratch dataspace.",
			},
			&cli.StringFlag{
				Name:    "url",
				Aliases: []string{"u"},
				Value:   "https://cloud.humio.com/",
				Usage:   "URL for the Humio server to stream to. `URL` must be a valid URL and end with slash (/).",
				EnvVars: []string{"HUMIO_URL"},
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "A name to make it easier to find results for this stream in your dataspace. e.g. @name=MyName\nIf `NAME` is not specified and you are tailing a file, the filename is used.",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "ingesttoken",
				Subcommands: []*cli.Command{
					{
						Name: "create",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Usage:   "Token name",
							},
							&cli.StringFlag{
								Name:    "parser",
								Aliases: []string{"p"},
								Usage:   "Parser name",
							},
						},
						Action: ingesttoken_create,
					},
				},
			},
			{
				Name: "parser",
				Subcommands: []*cli.Command{
					{
						Name: "create",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Usage:   "Parser name",
							},
							&cli.StringSliceFlag{
								Name:    "query",
								Aliases: []string{"q"},
								Usage:   "Query string",
							},
							&cli.BoolFlag{
								Name:    "force",
								Aliases: []string{"f"},
								Usage:   "Overwrite existing parser",
							},
						},
						Action: parser_create,
					},
				},
			},
			{
				Name: "ingest",
				Action: ingest,
			},
		},
	}

	app.Version = version

	app.Run(os.Args)
}

////////////////////////////////////////////////////////////////////////////////
///// Ingest command ///////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////


type config struct {
	ServerURL   string
	AuthToken   string
	DataspaceID string
	SessionID   string
	Name        string
}

func ingest(c *cli.Context) error {
	filepath := c.Args().Get(0)

	name := ""
	if c.String("name") == "" && filepath != "" {
		name = filepath
	} else {
		name = c.String("name")
	}

	u, _ := uuid.NewV4()

	sessionConfig := config{
		DataspaceID: c.String("dataspace"),
		AuthToken:   c.String("token"),
		ServerURL:   c.String("url"),
		Name:        name,
		SessionID:   u.String() ,
	}

	if sessionConfig.AuthToken == "" {
		log.Fatal("No AuthToken provided. See the -t option.")
	}

	// Open the browser (First so it has a chance to load)
	// key := ""
	// if name == "" {
	// 	key = "%40session%3D" + sessionConfig.SessionID
	// } else {
	// 	key = "%40name%3D" + sessionConfig.Name
	// }
	// open.Run(sessionConfig.ServerURL + sessionConfig.DataspaceID + "/search?live=true&start=1d&query=" + key)

	startSending(sessionConfig)

	if filepath != "" {
		tailFile(sessionConfig, filepath)
	} else {
		streamStdin(sessionConfig)
	}
	return nil
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

func sendBatch(config config, events []event) {

	lineJSON, marshalErr := json.Marshal([1]eventList{
		eventList{
			Tags:   map[string]string{},
			Events: events,
		}})

	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	ingestURL := config.ServerURL + "api/v1/dataspaces/" + config.DataspaceID + "/ingest"
	lineReq, reqErr := http.NewRequest("POST", ingestURL, bytes.NewBuffer(lineJSON))
	lineReq.Header.Set("Authorization", "Bearer "+config.AuthToken)
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

func startSending(config config) {
	go func() {
		var batch []event
		for {
			select {
			case v := <-events:
				batch = append(batch, v)
				if len(batch) >= batchLimit {
					sendBatch(config, batch)
					batch = batch[:0]
				}
			default:
				if len(batch) > 0 {
					sendBatch(config, batch)
					batch = batch[:0]
				}
				// Avoid busy waiting
				batch = append(batch, <-events)
			}
		}
	}()
}

func sendLine(config config, line string) {
	theEvent := event{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Attributes: map[string]string{
			"@session": config.SessionID,
			"@name":    config.Name,
		},
		RawString: line,
	}

	events <- theEvent
}

func tailFile(sessionConfig config, filepath string) {

	// Join Tail

	t, err := tail.TailFile(filepath, tail.Config{Follow: true})

	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		sendLine(sessionConfig, line.Text)
	}

	tailError := t.Wait()

	if tailError != nil {
		log.Fatal(tailError)
	}
}

func streamStdin(sessionConfig config) {
	log.Println("Humio Attached to StdIn")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sendLine(sessionConfig, scanner.Text())
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


////////////////////////////////////////////////////////////////////////////////
///// Ingesttoken command //////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func ingesttoken_create(c *cli.Context) error {
	name := c.String("name")
	if name == "" {
		panic("Missing name argument")
	}
	parser := c.String("parser")

	body := ""
	if parser == "" {
		body = `{"name": "`+name+`"}`
	}	else {
		body = `{"name": "`+name+`", "parser": "`+parser+`"}`
	}
	url := "http://localhost:8080/api/v1/repositories/developer/ingesttokens"

	resp, clientErr := postJson(url, body, "notoken")

	if clientErr != nil {
		panic(clientErr)
	}
	if resp.StatusCode >= 400 {
		responseData, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
		log.Print(resp.StatusCode)
		log.Fatal(string(responseData))
	}
	resp.Body.Close()

	return nil
}


////////////////////////////////////////////////////////////////////////////////
///// Parser create command ////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func parser_create(c *cli.Context) error {
	name := c.String("name")
	if name == "" {
		panic("Missing name argument")
	}
	query := c.StringSlice("query")
	log.Print("----->")
	log.Print(query)
	if len(query) != 1 {
		panic("Missing query argument")
	}
	log.Print(len(query))
	log.Print(query[0])

	body := `{"parser": "`+query[0]+`", "kind": "humio", "parseKeyValues": false, "dateTimeFields": ["@timestamp"]}`
	url := "http://localhost:8080/api/v1/repositories/developer/parsers/"+name

	resp, clientErr := postJson(url, body, "notoken")

	if clientErr != nil {
		panic(clientErr)
	}
	if resp.StatusCode == 409 && c.Bool("force") {
		resp, clientErr = putJson(url, body, "notoken")
	}

	if resp.StatusCode >= 400 {
		responseData, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
		log.Print(resp.StatusCode)
		log.Fatal(string(responseData))
	}
	log.Print(resp)
	resp.Body.Close()

	return nil
}


////////////////////////////////////////////////////////////////////////////////
///// Utils ////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func postJson(url string, jsonStr string, token string) (*http.Response, error) {
	log.Print(url)
	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}

func putJson(url string, jsonStr string, token string) (*http.Response, error) {
	log.Print(url)
	req, reqErr := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}