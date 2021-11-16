package main

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/viperkey"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
func newProfilesRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <profile>",
		Short: "Remove a configuration profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profileName := args[0]

			profiles := viper.GetStringMap(viperkey.Profiles)
			if profiles[profileName] == nil {
				cmd.Println("profile not found")
				os.Exit(0)
			}

			delete(profiles, profileName)
			err := saveConfig()
			exitOnError(cmd, err, "Error saving config")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed profile: %q\n", profileName)
		},
	}

	return cmd
}
