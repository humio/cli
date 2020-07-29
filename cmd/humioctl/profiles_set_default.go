package main

import (
	"fmt"
	"strconv"

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
			viper.Set("address", profile.address)
			viper.Set("token", profile.token)
			viper.Set("ca-certificate", string(profile.caCertificate))
			viper.Set("insecure", strconv.FormatBool(profile.insecure))

			saveErr := saveConfig()
			exitOnError(cmd, saveErr, "error saving config")

			out.Info(fmt.Sprintf("Default profile set to '%s'", profileName))

			cmd.Println()
			out.Output("Address: " + viper.GetString("address"))
			cmd.Println()
		},
	}

	return cmd
}

func loadProfile(profileName string) (*login, error) {
	profiles := viper.GetStringMap("profiles")
	profileData := profiles[profileName]

	if profileData == nil {
		return nil, fmt.Errorf("unknown profile %s", profileName)
	}

	insecure, err := strconv.ParseBool(getMapKey(profileData, "insecure"))
	if err != nil {
		return nil, err
	}

	profile := login{
		address:       getMapKey(profileData, "address"),
		token:         getMapKey(profileData, "token"),
		caCertificate: []byte(getMapKey(profileData, "ca-certificate")),
		insecure:      insecure,
	}

	return &profile, nil
}
