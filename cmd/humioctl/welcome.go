package main

import (
	"fmt"
	"os"

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
			profiles := viper.GetStringMap("profiles")
			out := prompt.NewPrompt(cmd.OutOrStdout())

			if profiles["default"] != nil && !out.Confirm("Your system is already set up for Humio. Re-initialize?") {
				os.Exit(0)
			}

			owl := "[purple]" + prompt.Owl() + "[reset]"
			out.Print((prompt.Colorize(owl)))
			out.Output("")
			out.Title("Welcome to Humio")
			out.Output("")
			out.Description("This will guide you through setting up the Humio CLI.")
			out.Output("")

			profile, err := collectProfileInfo(cmd)
			exitOnError(cmd, err, "failed to collect profile info")

			viper.Set("address", profile.address)
			viper.Set("token", profile.token)

			addAccount(out, "default", profile)

			configFile := viper.ConfigFileUsed()
			cmd.Println(prompt.Colorize("==> Writing settings to: [purple]" + configFile + "[reset]"))

			if saveErr := saveConfig(); saveErr != nil {
				cmd.Println(fmt.Errorf("error saving config: %s", saveErr))
				os.Exit(1)
			}

			cmd.Println()
			out.Description("The authentication info has been saved to the profile 'default'.")
			out.Description("If you work with multiple user accounts or Humio servers you can")
			out.Description("add more profiles using `humio profiles add <name>`.")

			cmd.Println()
			out.Output("Bye bye now! 🎉")
			cmd.Println()
		},
	}

	return cmd
}
