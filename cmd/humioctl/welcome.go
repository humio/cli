package main

import (
	"os"

	"github.com/humio/cli/internal/viperkey"
	"github.com/humio/cli/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
func newWelcomeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "welcome",
		Hidden: true,
		Short:  "Creates the 'default' profile",
		Args:   cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			profiles := viper.GetStringMap(viperkey.Profiles)
			out := prompt.NewPrompt(cmd.OutOrStdout())

			if profiles["default"] != nil && !out.Confirm("Your system is already set up for Humio. Re-initialize?") {
				os.Exit(0)
			}

			owl := "[purple]" + prompt.Owl() + "[reset]"
			out.Print(prompt.Colorize(owl))
			out.BlankLine()
			out.Title("Welcome to Humio")
			out.BlankLine()
			out.Description("This will guide you through setting up the Humio CLI.")
			out.BlankLine()

			profile, err := collectProfileInfo(cmd)
			exitOnError(cmd, err, "Failed to collect profile info")

			viper.Set(viperkey.Address, profile.address)
			viper.Set(viperkey.Token, profile.token)

			addAccount("default", profile)

			configFile := viper.ConfigFileUsed()
			out.Print(prompt.Colorize("==> Writing settings to: [purple]" + configFile + "[reset]"))

			err = saveConfig()
			exitOnError(cmd, err, "Error saving file")

			out.BlankLine()
			out.Description("The authentication info has been saved to the profile 'default'.")
			out.Description("If you work with multiple user accounts or Humio servers you can")
			out.Description("add more profiles using `humio profiles add <name>`.")

			out.BlankLine()
			out.Print("Bye bye now! 🎉")
			out.BlankLine()
		},
	}

	return cmd
}
