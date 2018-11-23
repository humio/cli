package command

import (
	"fmt"
	"strings"
)

type UsersListCommand struct {
	Meta
}

func (f *UsersListCommand) Help() string {
	helpText := `
Usage: humio users list

  Lists all users. This command requires root permissions on your access token.

  To see members in a repository or view use:

    $ humio members list <repo>

  General Options:

  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (f *UsersListCommand) Synopsis() string {
	return "List all user in the cluster."
}

func (f *UsersListCommand) Name() string { return "users list" }

func (f *UsersListCommand) Run(args []string) int {
	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check that we got one argument
	args = flags.Args()
	if l := len(args); l != 0 {
		f.Ui.Error("This command takes no arguments")
		f.Ui.Error(commandErrorText(f))
		return 1
	}

	// Get the HTTP client
	client, err := f.Meta.Client()
	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	users, err := client.Users().List()

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error fetching token list: %s", err))
		return 1
	}

	rows := make([]string, len(users))
	for i, user := range users {
		rows[i] = formatSimpleAccount(user)
	}

	printTable(append([]string{
		"Username | Name | Root | Created"},
		rows...,
	))

	return 0
}
