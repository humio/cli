package main

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/viperkey"
	"github.com/humio/cli/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type login struct {
	address       string
	token         string
	username      string
	caCertificate string
	insecure      bool
}

// usersCmd represents the users command
func newProfilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "List and manage configuration profiles.",
		Long: `If you interact with multiple Humio clusters or
use multiple user accounts you can save each profile for quickly switching between them.

If called without a subcommand this command will list saved profiles.

Adding a profile:

  $ humioctl profiles add <name>

You can change the default profile using:

  $ humioctl profiles set-default <name>
    `,
		Args: cobra.ExactArgs(0),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			profiles := viper.GetStringMap(viperkey.Profiles)

			for name, data := range profiles {
				login := mapToLogin(data)
				if isCurrentAccount(login.address, login.token) {
					cmd.Println(prompt.Colorize(fmt.Sprintf("* [purple]%s (%s) - %s[reset]", name, login.username, login.address)))
				} else {
					cmd.Println(fmt.Sprintf("  %s (%s) - %s", name, login.username, login.address))
				}
			}

			if len(profiles) == 0 {
				cmd.Println("You have no saved profiles")
				cmd.Println()
				cmd.Println("Use `humio profiles add <name>` to add one.")
			}

			cmd.Println()

			return nil, nil
		}),
	}

	cmd.AddCommand(newProfilesAddCmd())
	cmd.AddCommand(newProfilesRemoveCmd())
	cmd.AddCommand(newProfilesSetDefaultCmd())

	return cmd
}
