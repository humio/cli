package command

import (
	"fmt"
	"strings"
)

type TokensRemoveCommand struct {
	Meta
}

func (f *TokensRemoveCommand) Help() string {
	helpText := `
Usage: humio ingest-tokens rm <repo> <token name>

  Removes the ingest token with name <token name> from repository <repo>.

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *TokensRemoveCommand) Synopsis() string {
	return "Removes an ingest token from a repository."
}

func (f *TokensRemoveCommand) Name() string { return "ingest-tokens remove" }

func (f *TokensRemoveCommand) Run(args []string) int {

	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check that we got one argument
	args = flags.Args()
	if l := len(args); l != 2 {
		f.Ui.Error("This command takes two arguments: <repo> <token name>")
		f.Ui.Error(commandErrorText(f))
		return 1
	}

	repo := args[0]
	name := args[1]

	// Get the HTTP client
	client, err := f.Meta.Client()
	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	err = client.IngestTokens().Remove(repo, name)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error removing ingest token: %s", err))
		return 1
	}

	return 0
}
