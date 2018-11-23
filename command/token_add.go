package command

import (
	"fmt"
	"strings"

	"github.com/ryanuber/columnize"
)

type TokensAddCommand struct {
	Meta
}

func (f *TokensAddCommand) Help() string {
	helpText := `
Usage: humio ingest-tokens add <repo> <token name>

  Adds a new ingest token with name <token name> to a repository <repo>.

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *TokensAddCommand) Synopsis() string {
	return "Adds a new ingest token to a repository."
}

func (f *TokensAddCommand) Name() string { return "ingest-tokens add" }

func (f *TokensAddCommand) Run(args []string) int {
	var parserName stringPtrFlag

	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }
	flags.Var(&parserName, "parser", "Assign a parser to the token.")

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

	token, err := client.IngestTokens().Add(repo, name, parserName.value)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error adding ingest token: %s", err))
		return 1
	}

	var output []string
	output = append(output, "Name | Token | Assigned Parser")
	output = append(output, fmt.Sprintf("%v | %v | %v", token.Name, token.Token, valueOrEmpty(token.AssignedParser)))

	table := columnize.SimpleFormat(output)

	fmt.Println()
	fmt.Println(table)
	fmt.Println()

	return 0
}
