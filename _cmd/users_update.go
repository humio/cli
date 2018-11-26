package command

import (
	"fmt"
	"strings"

	"github.com/humio/cli/api"
)

type UsersUpdateCommand struct {
	Meta
}

func (f *UsersUpdateCommand) Help() string {
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

func (f *UsersUpdateCommand) Synopsis() string {
	return "Shows details about a user."
}

func (f *UsersUpdateCommand) Name() string { return "users update" }

func (f *UsersUpdateCommand) Run(args []string) int {
	var root boolPtrFlag

	flags := f.Meta.FlagSet(f.Name(), FlagSetClient)
	flags.Usage = func() { f.Ui.Output(f.Help()) }
	flags.Var(&root, "root", "If true grants root access to the user.")

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

	user, err := client.Users().Update(username, api.UserChangeSet{IsRoot: root.value})

	if err != nil {
		f.Ui.Error(fmt.Sprintf("Error updating user: %s", err))
		return 1
	}

	printUserTable(user)

	return 0
}
