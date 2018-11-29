package cmd

import (
	"fmt"

	"github.com/humio/cli/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
func newLoginListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "Lists saved accounts.",
		Long:  `Lists accounts shored in your ~/.humio/config.yaml file.`,
		Run: func(cmd *cobra.Command, args []string) {
			accounts := viper.GetStringMap("accounts")

			for name, data := range accounts {
				login := mapToLogin(data)
				if isCurrentAccount(login.address, login.token) {
					cmd.Println(prompt.Colorize(fmt.Sprintf("[purple]%s (%s) - %s[reset]", name, login.username, login.address)))
				} else {
					cmd.Println(fmt.Sprintf("%s (%s) - %s", name, login.username, login.address))
				}
			}

			if len(accounts) == 0 {
				cmd.Println("You have no saved accounts")
			}

			cmd.Println()
		},
	}

	return cmd
}

func isCurrentAccount(addr string, token string) bool {
	return viper.GetString("address") == addr && viper.GetString("token") == token
}
