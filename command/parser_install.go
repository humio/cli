package command

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/humio/cli/api"
	"gopkg.in/yaml.v2"
)

type ParsersInstallCommand struct {
	Meta
}

func (f *ParsersInstallCommand) Help() string {
	helpText := `
Usage: humio parsers rm <repo> <parser>

  Removes (uninstalls) the parser <parser> from <repo>.

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *ParsersInstallCommand) Synopsis() string {
	return "Remove a parser from a repository"
}

func (f *ParsersInstallCommand) Name() string { return "parsers rm" }

func (f *ParsersInstallCommand) Run(args []string) int {
	var content []byte
	var readErr error
	var verbose, force bool
	var filePath, url string

	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }
	flags.BoolVar(&verbose, "verbose", false, "")
	flags.BoolVar(&force, "force", false, "")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()

	// Check that we got the right number of argument
	if l := len(args); l == 1 {
		flags.StringVar(&filePath, "file", "", "")
		flags.StringVar(&url, "url", "", "")

		if filePath != "" {
			content, readErr = getParserFromFile(filePath)
		} else if url != "" {
			content, readErr = getUrlParser(url)
		} else {
			return 1
		}
	} else if l := len(args); l != 2 {
		f.Ui.Error("This command takes two arguments: <repo> <parser>")
		f.Ui.Error(commandErrorText(f))
		return 1
	} else {
		parserName := args[1]
		content, readErr = getGithubParser(parserName)
	}

	if readErr != nil {
		f.Ui.Error(fmt.Sprintf("Failed to get parser: %s", readErr))
		f.Ui.Error(commandErrorText(f))
	}

	parser := api.Parser{}
	err := yaml.Unmarshal(content, &parser)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("The parser's format was invalid: %s", readErr))
		f.Ui.Error(commandErrorText(f))
	}

	if flags.Lookup("name") != nil {
		parser.Name = *flags.String("name", "", "")
	}

	// Get the HTTP client
	client, err := f.Meta.Client()
	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	reposistoryName := args[0]

	err = client.Parsers().Add(reposistoryName, &parser, force)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error installing parser: %s", err))
		return 1
	}

	return 0
}

func getParserFromFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func getGithubParser(parserName string) ([]byte, error) {
	url := "https://raw.githubusercontent.com/humio/community/master/parsers/" + parserName + ".yaml"
	return getUrlParser(url)
}

func getUrlParser(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
