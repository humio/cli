package command

import (
	"fmt"
	"log"
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

func ensureRepo(server server) {
	if server.Repo == "" {
		log.Fatal("The command requires the repository to be specified.")
	}
}

func ensureURL(server server) {
	if server.URL == "" {
		log.Fatal("You must specify the URL of the Humio server.")
	}
}

func ensureToken(server server) {
	if server.Token == "" {
		exit("API Token not set.")
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
