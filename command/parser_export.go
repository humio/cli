package command

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type ParsersExportCommand struct {
	Meta
}

func (f *ParsersExportCommand) Help() string {
	helpText := `
Usage: humio parsers export <repo> <parser>

  Export a parser <parser> in <repo> to a file.

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *ParsersExportCommand) Synopsis() string {
	return "Exports a parser to a file"
}

func (f *ParsersExportCommand) Name() string { return "parsers export" }

func (f *ParsersExportCommand) Run(args []string) int {
	var outputName string

	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check that we got one argument
	args = flags.Args()
	if l := len(args); l != 2 {
		f.Ui.Error("This command takes two arguments: <repo> <parser>")
		f.Ui.Error(commandErrorText(f))
		return 1
	}

	repo := args[0]
	parserName := args[1]

	if flags.Lookup("output") != nil {
		flags.StringVar(&outputName, "output", "", "")
	} else {
		outputName = parserName
	}

	// Get the HTTP client
	client, err := f.Meta.Client()
	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	parser, err := client.Parsers().Get(repo, parserName)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error fetching parsers: %s", err))
		return 1
	}

	yamlData, yamlErr := yaml.Marshal(&parser)

	if yamlErr != nil {
		f.Ui.Error(fmt.Sprintf("Failed to serialize the parser: %s", yamlErr))
		return 1
	}

	outFilePath := outputName + ".yaml"

	err = ioutil.WriteFile(outFilePath, yamlData, 0644)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error saving the parser file: %s", err))
		return 1
	}

	return 0
}
