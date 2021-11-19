package humioctl

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/humio/cli/cmd/humioctl/internal/viperkey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newProfilesSetDefaultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-default <profile>",
		Short: "Choose one of your profiles to be the default.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profileName := args[0]

			profile, err := loadProfile(profileName)
			helpers.ExitOnError(cmd, err, "Profile not found")
			viper.Set(viperkey.Address, profile.address)
			viper.Set(viperkey.Token, profile.token)
			viper.Set(viperkey.CACertificateFile, profile.caCertificate)
			viper.Set(viperkey.Insecure, profile.insecure)

			err = saveConfig()
			helpers.ExitOnError(cmd, err, "Error saving config")

			fmt.Fprintf(cmd.OutOrStdout(), "Default profile set to %q\n", profileName)
			fmt.Fprintf(cmd.OutOrStdout(), "Address: %s\n", viper.GetString(viperkey.Address))
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

	insecureFromProfileData, _ := profileData[viperkey.Insecure].(bool) // false if not found in map, or type isn't bool

	profile := login{
		address:       getMapKeyString(profileData, viperkey.Address),
		token:         getMapKeyString(profileData, viperkey.Token),
		caCertificate: getMapKeyString(profileData, viperkey.CACertificate),
		insecure:      insecureFromProfileData,
	}

	return &profile, nil
}
