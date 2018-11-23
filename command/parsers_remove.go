package command

import (
	"fmt"
	"strings"
)

type ParsersRemoveCommand struct {
	Meta
}

func (f *ParsersRemoveCommand) Help() string {
	helpText := `
Usage: humio parsers rm <repo> <parser>

  Removes (uninstalls) the parser <parser> from <repo>.

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *ParsersRemoveCommand) Synopsis() string {
	return "Remove a parser from a repository"
}

func (f *ParsersRemoveCommand) Name() string { return "parsers rm" }

func (f *ParsersRemoveCommand) Run(args []string) int {

	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check that we got two argument
	args = flags.Args()
	if l := len(args); l != 2 {
		f.Ui.Error("This command takes two arguments: <repo> <parser>")
		f.Ui.Error(commandErrorText(f))
		return 1
	}

	repo := args[0]
	parser := args[1]

	// Get the HTTP client
	client, err := f.Meta.Client()
	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	err = client.Parsers().Remove(repo, parser)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error removing parser: %s", err))
		return 1
	}

	return 0
}
