package main

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/viperkey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
func newProfilesRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <profile-name> [flags]",
		Short: "Remove a configuration profile",
		Args:  cobra.ExactArgs(1),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			profileName := args[0]

			profiles := viper.GetStringMap(viperkey.Profiles)

			if profiles[profileName] == nil {
				return nil, humioErrorExitCode{fmt.Errorf("profile not found"), 0}
			}

			delete(profiles, profileName)

			saveErr := saveConfig()
			if saveErr != nil {
				return nil, fmt.Errorf("error saving config: %w", saveErr)
			}

			return "Profile removed: " + profileName, nil
		}),
	}

	return cmd
}
