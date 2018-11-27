package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/humio/cli/api"
	"github.com/humio/cli/prompt"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [flags]",
		Short: "Set up your environment to use Humio.",
		Long: `This command initializes your ~/.humio/config.yaml file with the information needed to interact with a Humio cluster.

You can also specify a differant config file to initialize by setting the --config flag.
In the config file already exists, the settings will be merged into the existing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var addr, token string
			var err error

			prompt.Output("")
			owl := "[purple]" + prompt.Owl() + "[reset]"
			fmt.Print((prompt.Colorize(owl)))
			prompt.Output("")
			prompt.Title("Welcome to Humio")
			prompt.Output("")
			prompt.Description("This will guide you through setting up the Humio CLI.")
			prompt.Output("")

			prompt.Info("Which Humio instance should we talk to?") // INFO
			prompt.Output("")
			prompt.Description("If you are not using Humio Cloud enter the address of your Humio installation,")
			prompt.Description("e.g. http://localhost:8080/ or https://humio.example.com/")

			for true {
				prompt.Output("")
				prompt.Output("Default: https://cloud.humio.com/ [Hit Enter]")
				addr, err = prompt.Ask("Humio Address")

				if addr == "" {
					addr = "https://cloud.humio.com/"
				}

				if err != nil {
					return fmt.Errorf("error reading humio server address: %s", err)
				}

				// Make sure it is a valid URL and that
				// we always end in a slash.
				_, urlErr := url.ParseRequestURI(addr)

				if urlErr != nil {
					prompt.Error("The valus must be a valid URL.")
					continue
				}

				if !strings.HasSuffix(addr, "/") {
					addr = addr + "/"
				}

				clientConfig := api.DefaultConfig()
				clientConfig.Address = addr
				client, apiErr := api.NewClient(clientConfig)

				if apiErr != nil {
					return (fmt.Errorf("error initializing the http client: %s", apiErr))
				}

				prompt.Output("")
				fmt.Print("==> Testing Connection...")

				status, statusErr := client.Status()

				if statusErr != nil {
					fmt.Println(prompt.Colorize("[[red]Failed[reset]]"))
					prompt.Output("")
					prompt.Error(fmt.Sprintf("Could not connect to the Humio server: %s\nIs the address connect and reachable?", statusErr))
					continue
				}

				if status.Status != "ok" {
					fmt.Println(prompt.Colorize("[[red]Failed[reset]]"))
					return (fmt.Errorf("The server reported that is is malfunctioning, status: %s", status.Status))
				} else {
					fmt.Println(prompt.Colorize("[[green]Ok[reset]]"))
				}

				fmt.Println("")
				break
			}

			prompt.Info("Paste in your Personal API Token")
			prompt.Description("")
			prompt.Description("To use Humio's CLI you will need to get a copy of your API Token.")
			prompt.Description("The API token can be found in your 'Account Settings' section of the UI.")
			prompt.Description("If you are running Humio without authorization just leave the API Token field empty.")
			prompt.Description("")

			if prompt.Confirm("Would you like us to open a browser on the account page?") {
				open.Start(fmt.Sprintf("%ssettings", addr))

				prompt.Description("")
				prompt.Description(fmt.Sprintf("If the browser did not open, you can manually visit:"))
				prompt.Description(fmt.Sprintf("%ssettings", addr))
				prompt.Description("")
			}

			for true {
				token, err = prompt.AskSecret("API Token")

				if err != nil {
					return (fmt.Errorf("error reading token: %s", err))
				}

				// Create a new API client with the token
				config := api.DefaultConfig()
				config.Address = addr
				config.Token = token
				client, clientErr := api.NewClient(config)

				if clientErr != nil {
					return fmt.Errorf("error initializing the http client: %s", clientErr)
				}

				username, apiErr := client.Viewer().Username()

				if apiErr != nil {
					prompt.Error("authorization failed, try another token")
					continue
				}

				fmt.Println("")
				fmt.Println("")
				fmt.Println(prompt.Colorize(fmt.Sprintf("==> Logged in as: [purple]%s[reset]", username)))
				fmt.Println("")
				break
			}

			viper.Set("address", addr)
			viper.Set("token", token)

			configFile := viper.ConfigFileUsed()

			fmt.Println(prompt.Colorize("==> Writing settings to: [purple]" + configFile + "[reset]"))

			if writeErr := viper.WriteConfig(); writeErr != nil {
				if os.IsNotExist(writeErr) {
					dirName := filepath.Dir(configFile)
					if dirErr := os.MkdirAll(dirName, 0700); dirErr != nil {
						return fmt.Errorf("error creating config directory %s: %s", dirName, dirErr)
					}
					if configFileErr := viper.WriteConfigAs(configFile); configFileErr != nil {
						return fmt.Errorf("error writing config file: %s", configFileErr)
					}
				}
			}

			fmt.Println("")
			prompt.Output("Bye bye now! 🎉")
			fmt.Println("")

			return nil
		},
	}

	return cmd
}
