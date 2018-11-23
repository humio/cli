package command

import (
	"fmt"
	"strings"

	"github.com/ryanuber/columnize"
)

type TokensListCommand struct {
	Meta
}

func (f *TokensListCommand) Help() string {
	helpText := `
Usage: humio tokens list <repo>

  List ingest tokens in a repository.

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *TokensListCommand) Synopsis() string {
	return "List ingest tokens in a repository."
}

func (f *TokensListCommand) Name() string { return "tokens list" }

func (f *TokensListCommand) Run(args []string) int {

	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check that we got one argument
	args = flags.Args()
	if l := len(args); l != 1 {
		f.Ui.Error("This command takes one arguments: <repo>")
		f.Ui.Error(commandErrorText(f))
		return 1
	}

	repo := args[0]

	// Get the HTTP client
	client, err := f.Meta.Client()
	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	tokens, err := client.IngestTokens().List(repo)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error fetching token list: %s", err))
		return 1
	}

	var output []string
	output = append(output, "Name | Token | Assigned Parser")
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		output = append(output, fmt.Sprintf("%v | %v | %v", token.Name, token.Token, valueOrEmpty(token.AssignedParser)))
	}

	table := columnize.SimpleFormat(output)

	fmt.Println()
	fmt.Println(table)
	fmt.Println()

	return 0
}

func valueOrEmpty(v string) string {
	if v == "" {
		return "-"
	}
	return v
}
