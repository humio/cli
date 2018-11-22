package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

type testCase struct {
	Input  string
	Output map[string]string
}

type parserConfig struct {
	Name        string
	Description string     `yaml:",omitempty"`
	Tests       []testCase `yaml:",omitempty"`
	Example     string     `yaml:",omitempty"`
	Script      string     `yaml:",flow"`
}

type ParsersCommand struct {
	Meta
}

func (f *ParsersCommand) Help() string {
	helpText := `
Usage: humio parsers <subcommand> [options] [args]
  This command groups subcommands for interacting with parsers.
  Users can install community parsers and save copies of parsers e.g. for
  keeping them version controlled.

  Installing a parser:
    $ humio parsers install <repo> <group/parser-name>

  You can find the list of parsers at: https://github.com/humio/community
  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(helpText)
}

func (f *ParsersCommand) Synopsis() string {
	return "Interact with parsers"
}

func (f *ParsersCommand) Name() string { return "parsers" }

func (f *ParsersCommand) Run(args []string) int {
	return cli.RunResultHelp
}
