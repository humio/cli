package main

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/viperkey"
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
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			profileName := args[0]
			out := prompt.NewPrompt(cmd.OutOrStdout())

			profile, loadErr := loadProfile(profileName)
			if loadErr != nil {
				return nil, fmt.Errorf("profile not found: %w", loadErr)
			}
			viper.Set(viperkey.Address, profile.address)
			viper.Set(viperkey.TokenFile, profile.token)
			viper.Set(viperkey.CACertificateFile, profile.caCertificate)
			viper.Set(viperkey.Insecure, strconv.FormatBool(profile.insecure))

			saveErr := saveConfig()
			if saveErr != nil {
				return nil, fmt.Errorf("error saving config: %w", saveErr)
			}

			out.Info(fmt.Sprintf("Default profile set to '%s'", profileName))

			cmd.Println()
			out.Output("Address: " + viper.GetString(viperkey.Address))
			cmd.Println()

			return nil, nil
		}),
	}

	return cmd
}

func loadProfile(profileName string) (*login, error) {
	profiles := viper.GetStringMap(viperkey.Profiles)
	profileData, ok := profiles[profileName].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unknown or invalid profile %s", profileName)
	}

	insecure, err := strconv.ParseBool(getMapKeyString(profileData, viperkey.Insecure))
	if err != nil {
		return nil, err
	}

	profile := login{
		address:       getMapKeyString(profileData, viperkey.Address),
		token:         getMapKeyString(profileData, viperkey.Token),
		caCertificate: getMapKeyString(profileData, viperkey.CACertificate),
		insecure:      insecure,
	}

	return &profile, nil
}
