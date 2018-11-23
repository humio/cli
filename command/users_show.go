package command

import (
	"fmt"
	"strings"

	"github.com/humio/cli/api"
)

type UsersShowCommand struct {
	Meta
}

func (f *UsersShowCommand) Help() string {
	helpText := `
Usage: humio users show <username>

  Shows details about a users. This command requires root access.

  To see members in a repository or view use:

    $ humio members show <repo> <username>

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *UsersShowCommand) Synopsis() string {
	return "Shows details about a user."
}

func (f *UsersShowCommand) Name() string { return "users show" }

func (f *UsersShowCommand) Run(args []string) int {
	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check that we got one argument
	args = flags.Args()
	if l := len(args); l != 1 {
		f.Ui.Error("This command takes one argument: <username>")
		f.Ui.Error(commandErrorText(f))
		return 1
	}

	username := args[0]

	// Get the HTTP client
	client, err := f.Meta.Client()
	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	user, err := client.Users().Get(username)

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error fetching token list: %s", err))
		return 1
	}

	printUserTable(user)

	return 0
}

func printUserTable(user api.User) {
	userData := []string{user.Username, user.FullName, user.CreatedAt, yesNo(user.IsRoot)}

	printTable([]string{
		"Username | Name | Created At | Is Root",
		strings.Join(userData, "|"),
	})
}
