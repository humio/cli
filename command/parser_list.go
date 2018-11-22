package command

import (
	"fmt"
	"strings"
)

func checkmark(v bool) string {
	if v {
		return "✓"
	}
	return ""
}

type ParsersListCommand struct {
	Meta
}

func (f *ParsersListCommand) Help() string {
	helpText := `
Usage: humio parsers list [options] <repository name>

  Lists the parsers installed in a repository.

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *ParsersListCommand) Synopsis() string {
	return "Interact with parsers"
}

func (f *ParsersListCommand) Name() string { return "parsers list" }

func (f *ParsersListCommand) Run(args []string) int {

	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check that we got one argument
	args = flags.Args()
	if l := len(args); l != 1 {
		f.Ui.Error("This command takes exactly one arguments")
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

	parsers, err := client.Parsers().List(repo)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error fetching parsers: %s", err))
		return 1
	}

	var output []string
	output = append(output, "Name | Custom")
	for i := 0; i < len(parsers); i++ {
		parser := parsers[i]
		output = append(output, fmt.Sprintf("%v | %v", parser.Name, checkmark(!parser.IsBuiltIn)))
	}

	printTable(output)

	return 0
}
