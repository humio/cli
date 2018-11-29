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

type Login struct {
	address  string
	token    string
	username string
}

// usersCmd represents the users command
func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [flags]",
		Short: "Set up your environment to use Humio.",
		Long: `This command initializes your ~/.humio/config.yaml file with the information needed to interact with a Humio cluster.

You can also specify a differant config file to initialize by setting the --config flag.
In the config file already exists, the settings will be merged into the existing.`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var addr, token string
			var err error
			accounts := viper.GetStringMap("accounts")

			if len(accounts) > 0 {
				prompt.Output("")
				prompt.Title("Add another account")
			} else {
				prompt.Output("")
				owl := "[purple]" + prompt.Owl() + "[reset]"
				fmt.Print((prompt.Colorize(owl)))
				prompt.Output("")
				prompt.Title("Welcome to Humio")
				prompt.Output("")
				prompt.Description("This will guide you through setting up the Humio CLI.")
				prompt.Output("")
			}

			if currentAddr := viper.GetString("address"); currentAddr != "" {
				currentToken := viper.GetString("token")

				if accounts != nil {
					var currentName string
					for name, data := range accounts {
						login := mapToLogin(data)
						if login.address == currentAddr && login.token == currentToken {
							currentName = name
						}
					}

					if currentName == "" {
						prompt.Output("")
						prompt.Output("You are currently logged in an unnamed account (address/token).")
						prompt.Output("Would you like to name the account before logging in with another?")
						prompt.Output("This way you can quickly switch between accounts and servers.")
						prompt.Output("")

						client := NewApiClient(cmd)
						username, usernameErr := client.Viewer().Username()

						prompt.Output("Server: " + currentAddr)

						if usernameErr != nil {
							prompt.Output(prompt.Colorize("Username: [red][Error] - Invalid token or bad network connection.[reset]"))
						} else {
							prompt.Output("Username: " + username)
						}

						prompt.Output("")
						if prompt.Confirm("Do you want to save this account?") {
							prompt.Output("")
							addAccount(cmd, currentAddr, currentToken, username)
						}
					}
				}
				prompt.Output("")
			}

			prompt.Info("Which Humio instance should we talk to?")
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

			var username string
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

				var apiErr error
				username, apiErr = client.Viewer().Username()

				if apiErr != nil {
					prompt.Error("Authentication failed, invalid token")

					if prompt.Confirm("Do you want to use another token?") {
						continue
					}
				}

				if username != "" {
					cmd.Println()
					cmd.Println()
					cmd.Println(prompt.Colorize(fmt.Sprintf("==> Logged in as: [purple]%s[reset]", username)))
				}

				cmd.Println()
				break
			}

			viper.Set("address", addr)
			viper.Set("token", token)

			if len(accounts) > 0 && prompt.Confirm("Would you like to give this account a name?") {
				addAccount(cmd, addr, token, username)
			}

			configFile := viper.ConfigFileUsed()
			cmd.Println(prompt.Colorize("==> Writing settings to: [purple]" + configFile + "[reset]"))

			if saveErr := saveConfig(); saveErr != nil {
				cmd.Println(fmt.Errorf("error saving config: %s", saveErr))
				os.Exit(1)
			}

			cmd.Println()
			prompt.Output("Bye bye now! 🎉")
			cmd.Println("")

			return nil
		},
	}

	cmd.AddCommand(newLoginListCmd())
	cmd.AddCommand(newLoginChooseCmd())

	return cmd
}

func saveConfig() error {
	configFile := viper.ConfigFileUsed()

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

	return nil
}

func addAccount(cmd *cobra.Command, address string, token string, username string) {
	prompt.Description("Example names: dev, prod, cloud, root")

	for true {
		newName, chooseErr := prompt.Ask("Choose a name for the account")
		if chooseErr != nil {
			cmd.Println(fmt.Errorf("error reading input %s", chooseErr))
			os.Exit(1)
		}
		newName = strings.TrimSpace(newName)
		if newName == "" {
			prompt.Error("Name cannot be blank.")
			continue
		} else if newName == "list" || newName == "add" || newName == "remove" {
			prompt.Error("The names `list`, `add`, and `remove` are reserved. Choose another.")
			continue
		}

		accounts := viper.GetStringMap("accounts")

		accounts[newName] = map[string]string{
			"address":  address,
			"token":    token,
			"username": username,
		}

		viper.Set("accounts", accounts)

		prompt.Output("")
		prompt.Info("Account Saved")
		prompt.Output("")
		break
	}
}

func mapToLogin(data interface{}) *Login {
	return &Login{
		address:  getMapKey(data, "address"),
		username: getMapKey(data, "username"),
		token:    getMapKey(data, "token"),
	}
}

func getMapKey(data interface{}, key string) string {
	m, ok1 := data.(map[string]interface{})
	if ok1 {
		v := m[key]
		vStr, ok2 := v.(string)

		if ok2 {
			return vStr
		}
	}

	return ""
}
