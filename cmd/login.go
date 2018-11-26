package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

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

			fmt.Println("This will guide you through setting up the Humio CLI.")
			fmt.Println("")

			for true {
				fmt.Println("1. Which Humio instance should we talk to?") // INFO
				fmt.Println("")
				fmt.Println("If you are not using Humio Cloud enter the address out your Humio installation,")
				fmt.Println("e.g. http://localhost:8080/ or https://humio.example.com/")
				fmt.Println("")
				fmt.Println("Default: https://cloud.humio.com/ [Hit Enter]")
				addr, err = prompt.Ask("Humio Address")

				if addr == "" {
					addr = "https://cloud.humio.com/"
				}

				if err != nil {
					return fmt.Errorf("error reading humio server address: %s", err)
				}

				// Make sure it is a valid URL and that
				// we always end in a slash.
				addrUrl, urlErr := url.Parse(addr)

				if urlErr != nil {
					fmt.Println("The valus must be a valid URL.") // ERROR
				}

				addr = fmt.Sprintf("%v", addrUrl)

				clientConfig := api.DefaultConfig()
				clientConfig.Address = addr
				client, apiErr := api.NewClient(clientConfig)

				if apiErr != nil {
					return (fmt.Errorf("error initializing the http client: %s", apiErr))
				}

				status, statusErr := client.Status()

				if statusErr != nil {
					fmt.Println(fmt.Errorf("Could not connect to the Humio server: %s\nIs the address connect and reachable?", statusErr)) // ERROR
					continue
				}

				if status.Status != "ok" {
					return (fmt.Errorf("The server reported that is is malfunctioning, status: %s", status.Status))
				}

				fmt.Println("")
				fmt.Println("==> Connection Successful") // INFO
				fmt.Println("")
				break
			}

			fmt.Println("")
			fmt.Println("2. Paste in your API Token") // INFO
			fmt.Println("")
			fmt.Println("To use Humio's CLI you will need to get a copy of your API Token.")
			fmt.Println("The API token can be found in your 'Account Settings' section of the UI.")
			fmt.Println("If you are running Humio without authorization just leave the API Token field empty.")
			fmt.Println("")
			prompt.Ask("We will now open the account page in a browser window. [Hit Any Key]")

			open.Start(fmt.Sprintf("%ssettings", addr))

			fmt.Println("")
			fmt.Println(fmt.Sprintf("If the browser did not open, you can manually visit:"))
			fmt.Println(fmt.Sprintf("%ssettings", addr))
			fmt.Println("")

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
					fmt.Println(fmt.Errorf("authorization failed, try another token")) // ERROR
					continue
				}

				fmt.Println("")
				fmt.Println(fmt.Sprintf("==> Login successful '%s' 🎉", username)) // INFO
				fmt.Println("")
				break
			}

			viper.Set("address", addr)
			viper.Set("token", token)

			// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
			// viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
			// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

			configFile := viper.ConfigFileUsed()

			fmt.Println("==> Writing settings to: " + configFile) // INFO

			if writeErr := viper.MergeInConfig(); writeErr != nil {
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
			fmt.Println("Bye bye now!")
			fmt.Println("")

			return nil
		},
	}

	return cmd
}
