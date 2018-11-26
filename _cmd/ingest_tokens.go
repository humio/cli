package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

type IngestTokenCommand struct {
	Meta
}

func (f *IngestTokenCommand) Help() string {
	helpText := `
Usage: humio ingest-tokens <subcommand> [options] [args]
  This command groups subcommands for interacting with ingest tokens
  for a repository.

  To create an ingest-token:

    $ humio ingest-token add -parser=<parser> <repo> <token name>

  where <parser> is the name of one of the parsers on <repo>.

  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(helpText)
}

func (f *IngestTokenCommand) Synopsis() string {
	return "Interact with ingest tokens"
}

func (f *IngestTokenCommand) Name() string { return "ingest-tokens" }

func (f *IngestTokenCommand) Run(args []string) int {
	return cli.RunResultHelp
}
