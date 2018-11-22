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
				Usage:   "The repository to interact with.",
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
				Name: "users",
				Subcommands: []*cli.Command{
					{
						Name:   "show",
						Action: command.UsersShow,
					},
					{
						Name: "update",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "root",
								Usage: "Grant root permission to the user.",
							},
						},
						Action: command.UpdateUser,
					},
					{
						Name:   "list",
						Action: command.UsersList,
					},
				},
			},
			{
				Name: "tokens",
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
				Name: "parsers",
				Subcommands: []*cli.Command{
					{
						Name:   "get",
						Action: command.ParserGet,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "output",
								Aliases: []string{"out", "o"},
								Usage:   "The file path where the parser file should be stored.",
							},
						},
					},
					{
						Name: "push",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Usage:   "Override the name in the parser file.",
							},
							&cli.BoolFlag{
								Name:    "update",
								Aliases: []string{"u"},
								Usage:   "If a parser exists with the same name update it.",
							},
						},
						Action: command.ParserPush,
					},
					{
						Name:   "remove",
						Action: command.ParserRemove,
					},
					{
						Name:   "list",
						Action: command.ParserList,
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
					&cli.StringFlag{
						Name:    "parser",
						Aliases: []string{"p"},
						Usage:   "The name of the parser to use for ingest. This will have no effect if you have assigned parser to the ingest token used.",
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
// Utils
////////////////////////////////////////////////////////////////////////////////

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
