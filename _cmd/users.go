package command

import (
	"strings"

	"github.com/humio/cli/api"
	"github.com/mitchellh/cli"
)

func formatSimpleAccount(account api.User) string {
	columns := []string{account.Username, account.FullName, yesNo(account.IsRoot), account.CreatedAt}
	return strings.Join(columns, " | ")
}

type UsersCommand struct {
	Meta
}

func (f *UsersCommand) Help() string {
	helpText := `
Usage: humio users <subcommand> [options] [args]
  This command groups subcommands for interacting with users and
  requires root permissions in the Humio cluster.

  To grant root access to another user:

    $ humio users update -root=true <username>

  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(helpText)
}

func (f *UsersCommand) Synopsis() string {
	return "Interact with users"
}

func (f *UsersCommand) Name() string { return "users" }

func (f *UsersCommand) Run(args []string) int {
	return cli.RunResultHelp
}

func printUserTable(user api.User) {
	userData := []string{user.Username, user.FullName, user.CreatedAt, yesNo(user.IsRoot)}

	printTable([]string{
		"Username | Name | Created At | Is Root",
		strings.Join(userData, "|"),
	})
}
