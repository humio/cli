package main

import (
	"fmt"
	"net/http"
	"os"
	"os/user"

	// "github.com/skratchdot/open-golang/open"
	"github.com/humio/cli/command"
	"github.com/joho/godotenv"
	"gopkg.in/urfave/cli.v2"
)

////////////////////////////////////////////////////////////////////////////////
///// Globals //////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// This is set by the release script.
var version = "master"
var client = &http.Client{}

type server struct {
	URL   string
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
				Usage:   "URL for the Humio server. `URL` must be a valid URL and end with slash (/).",
				EnvVars: []string{"HUMIO_URL"},
			},
		},
		Commands: []*cli.Command{
			{
				Name: "token",
				Subcommands: []*cli.Command{
					{
						Name: "add",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Usage:   "The name to assign to the token.",
							},
							&cli.StringFlag{
								Name:    "parser",
								Aliases: []string{"p"},
								Usage:   "The name of the parser to assign to the ingest token.",
							},
						},
						Action: command.TokenAdd,
					},
					{
						Name:   "list",
						Action: command.TokenList,
					},
				},
			},
			{
				Name: "parser",
				Subcommands: []*cli.Command{
					{
						Name: "push",
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
								Name:  "query-file",
								Usage: "File containing the query",
							},
							&cli.BoolFlag{
								Name:    "force",
								Aliases: []string{"f"},
								Usage:   "Overwrite existing parser",
							},
						},
						Action: command.ParserPush,
					},
				},
			},
			{
				Name:   "ingest",
				Action: command.Ingest,
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
	godotenv.Load(user.HomeDir + "/.humio-cli.env")
}

////////////////////////////////////////////////////////////////////////////////
///// Utils ////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
