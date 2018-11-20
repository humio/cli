package command

import (
	"fmt"
	"log"
	"net/http"
	"os"

	cli "gopkg.in/urfave/cli.v2"
)

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

func newPostRequest(config server, path string) {

	req, err := http.NewRequest("POST", config.URL+path, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+config.Token)
}

func ensureRepo(server server) {
	if server.Repo == "" {
		exit("Missing repository argument")
	}
}

func ensureURL(server server) {
	if server.URL == "" {
		exit("Missing url argument")
	}
}

func ensureToken(server server) {
	if server.Token == "" {
		exit("Missing API token argument")
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
