package cmd

import (
	"fmt"
	"os"

	"github.com/humio/cli/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newProfilesSetDefaultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-default <profile-name>",
		Short: "Choose one of your profiles to be the default.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profiles := viper.GetStringMap("profiles")
			profileName := args[0]
			out := prompt.NewPrompt(cmd.OutOrStdout())

			profile := profiles[profileName]

			if profile == nil {
				cmd.Println(fmt.Errorf("unknown profile %s", profileName))
				os.Exit(1)
			}

			address := getMapKey(profile, "address")
			token := getMapKey(profile, "token")
			viper.Set("address", address)
			viper.Set("token", token)

			saveErr := saveConfig()
			exitOnError(cmd, saveErr, "error saving config")

			out.Info(fmt.Sprintf("Default profile set to '%s'", profileName))

			cmd.Println()
			out.Output("Address: " + address)
			cmd.Println()
		},
	}

	return cmd
}
