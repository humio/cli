package main

import (
	"os"

	"github.com/humio/cli/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
func newProfilesRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <profile-name> [flags]",
		Short: "Remove a configuration profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profileName := args[0]

			out := prompt.NewPrompt(cmd.OutOrStdout())

			profiles := viper.GetStringMap("profiles")

			if profiles[profileName] == nil {
				cmd.Println("profile not found")
				os.Exit(0)
			}

			delete(profiles, profileName)

			saveErr := saveConfig()
			exitOnError(cmd, saveErr, "error saving config")

			out.Output("Profile removed: ", profileName)
		},
	}

	return cmd
}
