package cmd

import (
	"fmt"
	"os"

	"github.com/humio/cli/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
func newLoginChooseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "choose <account-name>",
		Short: "Choose one of your saved accounts to be the active one.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			accounts := viper.GetStringMap("accounts")
			accountName := args[0]

			account := accounts[accountName]

			if account == nil {
				cmd.Println(fmt.Errorf("unknown account %s", accountName))
				os.Exit(1)
			}

			address := getMapKey(account, "address")
			token := getMapKey(account, "token")
			username := getMapKey(account, "username")

			viper.Set("address", address)
			viper.Set("token", token)

			if saveErr := saveConfig(); saveErr != nil {
				fmt.Println(fmt.Errorf("error saving config: %s", saveErr))
				os.Exit(1)
			}
			prompt.Info(fmt.Sprintf("Switched to account: '%s'", accountName))
			cmd.Println()
			prompt.Output("Address: " + address)
			if username != "" {
				prompt.Output("Username: " + username)
			}
			cmd.Println()
		},
	}

	return cmd
}
