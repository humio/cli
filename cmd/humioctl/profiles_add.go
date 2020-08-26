package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/humio/cli/api"
	"github.com/humio/cli/prompt"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// usersCmd represents the users command
func newProfilesAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <profile-name> [flags]",
		Short: "Add a configuration profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			out := prompt.NewPrompt(cmd.OutOrStdout())

			profileName := args[0]

			profile, profileErr := collectProfileInfo(cmd)
			exitOnError(cmd, profileErr, "failed to collect profile info")

			addAccount(out, profileName, profile)

			saveErr := saveConfig()
			exitOnError(cmd, saveErr, "error saving config")
		},
	}

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

func addAccount(out *prompt.Prompt, newName string, profile *login) {
	profiles := viper.GetStringMap("profiles")

	profiles[newName] = map[string]string{
		"address":        profile.address,
		"token":          profile.token,
		"username":       profile.username,
		"ca_certificate": string(profile.caCertificate),
		"insecure":       strconv.FormatBool(profile.insecure),
	}

	viper.Set("profiles", profiles)
}

func mapToLogin(data interface{}) (*login, error) {
	insecure, err := strconv.ParseBool(getMapKey(data, "insecure"))
	if err != nil {
		return nil, err
	}
	return &login{
		address:       getMapKey(data, "address"),
		username:      getMapKey(data, "username"),
		token:         getMapKey(data, "token"),
		caCertificate: []byte(getMapKey(data, "caCertificate")),
		insecure:      insecure,
	}, nil
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

func collectProfileInfo(cmd *cobra.Command) (*login, error) {
	var addr, token, username, caCertificate string
	var insecure bool

	out := prompt.NewPrompt(cmd.OutOrStdout())
	out.Info("Which Humio instance should we talk to?")
	out.Output()
	out.Description("If you are not using Humio Cloud enter the address of your Humio installation,")
	out.Description("e.g. http://localhost:8080/ or https://humio.example.com/")

	for {
		var err error
		out.Output("")
		addr, err = out.Ask("Address (default: https://cloud.humio.com/ [Hit Enter])")
		exitOnError(cmd, err, "error reading humio server address")

		if addr == "" {
			addr = "https://cloud.humio.com/"
		}

		// Make sure it is a valid URL and that
		// we always end in a slash.
		_, urlErr := url.ParseRequestURI(addr)

		if urlErr != nil {
			out.Error("The value must be a valid URL.")
			continue
		}

		if !strings.HasSuffix(addr, "/") {
			addr = addr + "/"
		}

		clientConfig := api.DefaultConfig()
		clientConfig.Address = addr
		client, apiErr := api.NewClient(clientConfig)
		exitOnError(cmd, apiErr, "error initializing the API client")

		out.Output()
		cmd.Print("==> Testing Connection...")
		out.Output()
		_, statusErr := client.Status()
		if statusErr != nil {
			if strings.Contains(statusErr.Error(), "x509: certificate signed by unknown authority") {
				cmd.Println(prompt.Colorize("[[red]Failed[reset]] Certificate not signed by a trusted Certificate Authority."))
				out.Info("What is the absolute path to the CA certificate that should be used for TLS certificate validation?")
				out.Output()
				out.Description("If you require a custom CA certificate for validating the TLS certificate of the Humio cluster,")
				out.Description("specify the path to the file containing the CA certificate in PEM format.")
				out.Description("If left empty it is not possible to validate TLS certificate chain.")
				out.Output()

				caCertificateFilePath, err := out.Ask("Absolute path on local disk to CA certificate in PEM format")
				exitOnError(cmd, err, "error reading Humio CA certificate file path")
				if caCertificateFilePath != "" {
					// Read the file
					caCertContent, err := ioutil.ReadFile(caCertificateFilePath)
					exitOnError(cmd, err, "error reading Humio CA certificate file path")
					caCertificate = string(caCertContent)
					clientConfig.CACertificate = []byte(caCertificate)
					client, err = api.NewClient(clientConfig)
					exitOnError(cmd, err, "error initializing the API client")
				}
			}
		}

		out.Output()
		cmd.Print("==> Testing Connection...")
		out.Output()
		_, statusErr = client.Status()
		if statusErr != nil {
			if strings.Contains(statusErr.Error(), "x509: certificate is valid for") {
				cmd.Printf("%s: %s\n", prompt.Colorize("[[red]Failed[reset]] Certificate not valid for"), clientConfig.Address)
				out.Output()
				out.Info("Disable hostname verification for TLS connections?")
				out.Output()
				out.Description("By default all connections will verify the hostname, but this option allows you to disable this if required.")
				out.Output()
				insecureString, err := out.Ask("Do you want to disable hostname verification? Type 'yes' to disable hostname verification")
				exitOnError(cmd, err, "error reading humio ca certificate file path")
				if insecureString == "yes" {
					out.Output("Disabling hostname verification.")
					insecure = true
					clientConfig.Insecure = true
					client, err = api.NewClient(clientConfig)
					exitOnError(cmd, err, "error initializing the API client")
				}
			}
		}

		out.Output()
		cmd.Print("==> Testing Connection...")
		out.Output()
		status, statusErr := client.Status()

		if statusErr != nil {
			cmd.Println(prompt.Colorize("[[red]Failed[reset]]"))
			out.Output()
			out.Error(fmt.Sprintf("Could not connect to the Humio server: %s\nIs the address connect and reachable?", statusErr))
			continue
		}

		if status.IsDown() {
			cmd.Println(prompt.Colorize("[[red]Failed[reset]]"))
			cmd.Println(fmt.Errorf("The server reported that it is malfunctioning, status: %s", status.Status))
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
	out.Description("The API token can be found on the 'Account Settings' page of the UI.")
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

	for {
		var err error
		token, err = out.AskSecret("API Token")
		exitOnError(cmd, err, "error reading token")

		// Create a new API client with the token
		config := api.DefaultConfig()
		config.Address = addr
		config.Token = token
		config.CACertificate = []byte(caCertificate)
		config.Insecure = insecure

		client, clientErr := api.NewClient(config)

		exitOnError(cmd, clientErr, "error initializing the http client")

		var apiErr error
		username, apiErr = client.Viewer().Username()

		if apiErr != nil {
			out.Error(fmt.Sprintf("Authentication failed, invalid token: %s", apiErr))

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

	return &login{address: addr, token: token, username: username, caCertificate: []byte(caCertificate), insecure: insecure}, nil
}

func isCurrentAccount(addr string, token string) bool {
	return viper.GetString("address") == addr && viper.GetString("token") == token
}
