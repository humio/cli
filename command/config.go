package command

import (
	"fmt"
	"os"
)

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
