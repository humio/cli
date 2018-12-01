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
		Run: func(cmd *cobra.Command, args []string) {
			var addr, token string
			var err error
			accounts := viper.GetStringMap("accounts")
			out := prompt.NewPrompt(cmd.OutOrStdout())

			if len(accounts) > 0 {
				out.Output("")
				out.Title("Add another account")
			} else {
				out.Output("")
				owl := "[purple]" + prompt.Owl() + "[reset]"
				out.Print((prompt.Colorize(owl)))
				out.Output("")
				out.Title("Welcome to Humio")
				out.Output("")
				out.Description("This will guide you through setting up the Humio CLI.")
				out.Output("")
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
						out.Output()
						out.Output("You are currently logged in an unnamed account (address/token).")
						out.Output("Would you like to name the account before logging in with another?")
						out.Output("This way you can quickly switch between accounts and servers.")
						out.Output()

						client := NewApiClient(cmd)
						username, usernameErr := client.Viewer().Username()

						out.Output("Server: " + currentAddr)

						if usernameErr != nil {
							out.Output(prompt.Colorize("Username: [red][Error] - Invalid token or bad network connection.[reset]"))
						} else {
							out.Output("Username: " + username)
						}

						out.Output()
						if out.Confirm("Do you want to save this account?") {
							out.Output()
							addAccount(out, currentAddr, currentToken, username)
						}
					}
				}
				out.Output()
			}

			out.Info("Which Humio instance should we talk to?")
			out.Output()
			out.Description("If you are not using Humio Cloud enter the address of your Humio installation,")
			out.Description("e.g. http://localhost:8080/ or https://humio.example.com/")

			for true {
				out.Output("")
				addr, err = out.Ask("Address (default: https://cloud.humio.com/ [Hit Enter])")

				if addr == "" {
					addr = "https://cloud.humio.com/"
				}

				exitOnError(cmd, err, "error reading humio server address")

				// Make sure it is a valid URL and that
				// we always end in a slash.
				_, urlErr := url.ParseRequestURI(addr)

				if urlErr != nil {
					out.Error("The valus must be a valid URL.")
					continue
				}

				if !strings.HasSuffix(addr, "/") {
					addr = addr + "/"
				}

				clientConfig := api.DefaultConfig()
				clientConfig.Address = addr
				client, apiErr := api.NewClient(clientConfig)
				exitOnError(cmd, apiErr, "error initializing the API client")

				out.Output("")
				cmd.Print("==> Testing Connection...")

				status, statusErr := client.Status()

				if statusErr != nil {
					cmd.Println(prompt.Colorize("[[red]Failed[reset]]"))
					out.Output()
					out.Error(fmt.Sprintf("Could not connect to the Humio server: %s\nIs the address connect and reachable?", statusErr))
					continue
				}

				if status.Status != "ok" {
					cmd.Println(prompt.Colorize("[[red]Failed[reset]]"))
					cmd.Println(fmt.Errorf("The server reported that is is malfunctioning, status: %s", status.Status))
					os.Exit(1)
				} else {
					cmd.Println(prompt.Colorize("[[green]Ok[reset]]"))
				}

				fmt.Println("")
				break
			}

			out.Info("Paste in your Personal API Token")
			out.Output()
			out.Description("To use Humio's CLI you will need to get a copy of your API Token.")
			out.Description("The API token can be found in your 'Account Settings' section of the UI.")
			out.Description("If you are running Humio without authorization just leave the API Token field empty.")
			out.Output()

			if out.Confirm("Would you like us to open a browser on the account page?") {
				open.Start(fmt.Sprintf("%ssettings", addr))

				out.Output()
				out.Description(fmt.Sprintf("If the browser did not open, you can manually visit:"))
				out.Description(fmt.Sprintf("%ssettings", addr))
				out.Output()
			}

			out.Output()

			var username string
			for true {
				token, err = out.AskSecret("API Token")
				exitOnError(cmd, err, "error reading token")

				// Create a new API client with the token
				config := api.DefaultConfig()
				config.Address = addr
				config.Token = token
				client, clientErr := api.NewClient(config)

				exitOnError(cmd, clientErr, "error initializing the http client")

				var apiErr error
				username, apiErr = client.Viewer().Username()

				if apiErr != nil {
					out.Error("Authentication failed, invalid token")

					if out.Confirm("Do you want to use another token?") {
						continue
					}
				}

				if username != "" {
					out.Output()
					out.Output()
					cmd.Println(prompt.Colorize(fmt.Sprintf("==> Logged in as: [purple]%s[reset]", username)))
				}

				cmd.Println()
				break
			}

			viper.Set("address", addr)
			viper.Set("token", token)

			if len(accounts) > 0 && out.Confirm("Would you like to give this account a name?") {
				out.Output()
				addAccount(out, addr, token, username)
			}

			configFile := viper.ConfigFileUsed()
			cmd.Println(prompt.Colorize("==> Writing settings to: [purple]" + configFile + "[reset]"))

			if saveErr := saveConfig(); saveErr != nil {
				cmd.Println(fmt.Errorf("error saving config: %s", saveErr))
				os.Exit(1)
			}

			cmd.Println()
			out.Output("Bye bye now! 🎉")
			cmd.Println()
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

func addAccount(out *prompt.Prompt, address string, token string, username string) {
	out.Description("Example names: dev, prod, cloud, root")

	for true {
		newName, chooseErr := out.Ask("Choose a name for the account")
		if chooseErr != nil {
			out.Output(fmt.Errorf("error reading input %s", chooseErr))
			os.Exit(1)
		}
		newName = strings.TrimSpace(newName)
		if newName == "" {
			out.Error("Name cannot be blank.")
			continue
		} else if newName == "list" || newName == "add" || newName == "remove" {
			out.Error("The names `list`, `add`, and `remove` are reserved. Choose another.")
			continue
		}

		accounts := viper.GetStringMap("accounts")

		accounts[newName] = map[string]string{
			"address":  address,
			"token":    token,
			"username": username,
		}

		viper.Set("accounts", accounts)

		out.Output("")
		out.Info("Account Saved")
		out.Output("")
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
