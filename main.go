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
	"os/user"
	"os/signal"
	"syscall"
	"time"
	"strconv"

	"github.com/hpcloud/tail"
	"github.com/satori/go.uuid"
	// "github.com/skratchdot/open-golang/open"
	"gopkg.in/urfave/cli.v2"
	"github.com/joho/godotenv"
)


////////////////////////////////////////////////////////////////////////////////
///// Globals //////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// This is set by the release script.
var version = "master"
var client = &http.Client{}

type server struct {
	Url   string
	Token string
	Repo  string
}


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
				Usage:   "Your Humio API Token",
				EnvVars: []string{"HUMIO_API_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "repo",
				Aliases: []string{"r"},
				Usage:   "The repository to stream to.",
				EnvVars: []string{"HUMIO_REPO"},
			},
			&cli.StringFlag{
				Name:    "url",
				Usage:   "URL for the Humio server to stream to. `URL` must be a valid URL and end with slash (/).",
				EnvVars: []string{"HUMIO_URL"},
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
//					{
//						Name: "list",
//						Action: ingesttoken_list,
//					},
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
							&cli.StringSliceFlag{
								Name:    "query-file",
								Usage:   "File containing the query",
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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "A name to make it easier to find results for this stream in your repository. e.g. @name=MyName\nIf `NAME` is not specified and you are tailing a file, the filename is used.",
					},
				},
			},
		},
	}

	app.Version = version
	loadEnvFile()
	app.Run(os.Args)
}

func loadEnvFile() {
	user, userErr := user.Current()
	if userErr != nil {
		panic(userErr)
	}
	// Load the env file if it exists
	godotenv.Load(user.HomeDir+"/.humio-cli.env")
}


////////////////////////////////////////////////////////////////////////////////
///// Ingest command ///////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func ingest(c *cli.Context) error {
	filepath := c.Args().Get(0)

	name := ""
	if c.String("name") == "" && filepath != "" {
		name = filepath
	} else {
		name = c.String("name")
	}

	u, _ := uuid.NewV4()

	sessionID := u.String()

	server, _ := getServerConfig(c)

	ensureToken(server)

	// Open the browser (First so it has a chance to load)
	// key := ""
	// if name == "" {
	// 	key = "%40session%3D" + server.SessionID
	// } else {
	// 	key = "%40name%3D" + server.Name
	// }
	// open.Run(server.ServerURL + server.RepoID + "/search?live=true&start=1d&query=" + key)

	startSending(server)

	if filepath != "" {
		tailFile(server, name, sessionID, filepath)
	} else {
		streamStdin(server, name, sessionID)
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

func sendBatch(server server, events []event) {

	lineJSON, marshalErr := json.Marshal([1]eventList{
		eventList{
			Tags:   map[string]string{},
			Events: events,
		}})

	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	ingestURL := server.Url + "/api/v1/dataspaces/" + server.Repo + "/ingest"
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


////////////////////////////////////////////////////////////////////////////////
///// Ingesttoken create command ///////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func ingesttoken_create(c *cli.Context) error {
	server, _ := getServerConfig(c)

	ensureToken(server)
	ensureRepo(server)
	ensureUrl(server)

	name := c.String("name")
	if name == "" {
		exit("Missing name argument")
	}
	parser := c.String("parser")

	body := ""
	if parser == "" {
		body = `{"name": "`+name+`"}`
	}	else {
		body = `{"name": "`+name+`", "parser": "`+parser+`"}`
	}

	url := server.Url+"/api/v1/repositories/"+server.Repo+"/ingesttokens"

	resp, clientErr := postJson(url, body, server.Token)

	if clientErr != nil {
		panic(clientErr)
	}
	if resp.StatusCode >= 400 {
		_, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
//		fmt.Println(resp.StatusCode)
//		fmt.Println(string(responseData))
	}
	resp.Body.Close()

	return nil
}


////////////////////////////////////////////////////////////////////////////////
///// Ingesttoken list command /////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func ingesttoken_list(c *cli.Context) error {
	server, _ := getServerConfig(c)

	ensureToken(server)
	ensureRepo(server)
	ensureUrl(server)

	url := server.Url+"/api/v1/repositories/"+server.Repo+"/ingesttokens"

	resp, clientErr := getReq(url, server.Token)
	fmt.Println(resp.Body)
	if clientErr != nil {
		panic(clientErr)
	}
	if resp.StatusCode >= 400 {
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
		fmt.Println(body)
	}
	resp.Body.Close()

	return nil
}


////////////////////////////////////////////////////////////////////////////////
///// Parser create command ////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func parser_create(c *cli.Context) error {
	server, _ := getServerConfig(c)

	ensureToken(server)
	ensureRepo(server)
	ensureUrl(server)

	name := c.String("name")
	if name == "" {
		panic("Missing name argument")
	}

	query := ""

	fileNameSlices := c.StringSlice("query-file")
	if len(fileNameSlices) != 1 {
		querySlices := c.StringSlice("query")
		if len(querySlices) != 1 {
			panic("Missing query argument")
		} else {
			query = strconv.Quote(querySlices[0])
		}
	} else {
		file, readErr := ioutil.ReadFile(fileNameSlices[0])
		if readErr != nil {
			exit("Could not read file: "+fileNameSlices[0])
		}
		query = strconv.Quote(string(file))
	}

	body := `{"parser": `+query+`, "kind": "humio", "parseKeyValues": false, "dateTimeFields": ["@timestamp"]}`
	url := server.Url+"/api/v1/repositories/"+server.Repo+"/parsers/"+name
	resp, clientErr := postJson(url, body, server.Token)

	if clientErr != nil {
		panic(clientErr)
	}
	if resp.StatusCode == 409 && c.Bool("force") {
		resp, clientErr = putJson(url, body, server.Token)
	}

	if resp.StatusCode >= 400 {
		responseData, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			panic(readErr)
		}
		fmt.Println("Error: status code =", resp.StatusCode)
		fmt.Println(string(responseData))
	}
	//fmt.Println(resp)
	resp.Body.Close()

	return nil
}


////////////////////////////////////////////////////////////////////////////////
///// Utils ////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func postJson(url string, jsonStr string, token string) (*http.Response, error) {
	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}

func putJson(url string, jsonStr string, token string) (*http.Response, error) {
	req, reqErr := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}

func getReq(url string, token string) (*http.Response, error) {
	req, reqErr := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	fmt.Println(req)
	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}

func getServerConfig(c *cli.Context) (server, error) {
	server := server{
		Repo:  c.String("repo"),
		Token: c.String("token"),
		Url:   c.String("url"),
	}
	return server, nil
}

func ensureRepo(server server){
	if server.Repo == "" {
		exit("Missing repository argument")
	}
}

func ensureUrl(server server){
	if server.Url == "" {
		exit("Missing url argument")
	}
}

func ensureToken(server server){
	if server.Token == "" {
		exit("Missing API token argument")
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}