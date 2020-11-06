package main

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/viperkey"
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
			profileName := args[0]
			out := prompt.NewPrompt(cmd.OutOrStdout())

			profile, loadErr := loadProfile(profileName)
			exitOnError(cmd, loadErr, "profile not found")
			viper.Set(viperkey.Address, profile.address)
			viper.Set(viperkey.Token, profile.token)
			viper.Set(viperkey.CACertificateFile, profile.caCertificate)
			viper.Set(viperkey.Insecure, profile.insecure)

			saveErr := saveConfig()
			exitOnError(cmd, saveErr, "error saving config")

			out.Info(fmt.Sprintf("Default profile set to '%s'", profileName))

			cmd.Println()
			out.Output("Address: " + viper.GetString(viperkey.Address))
			cmd.Println()
		},
	}

	return cmd
}

func loadProfile(profileName string) (*login, error) {
	profiles := viper.GetStringMap(viperkey.Profiles)
	profileData, ok := profiles[profileName].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unknown or invalid profile %s", profileName)
	}

	insecure, _ := profileData[viperkey.Insecure].(bool) // false if not found in map, or type isn't bool

	profile := login{
		address:       getMapKeyString(profileData, viperkey.Address),
		token:         getMapKeyString(profileData, viperkey.Token),
		caCertificate: getMapKeyString(profileData, viperkey.CACertificate),
		insecure:      insecure,
	}

	return &profile, nil
}
