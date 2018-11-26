package command

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/humio/cli/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
)

type LoginCommand struct {
	Meta
}

func (f *LoginCommand) Help() string {
	helpText := `
Usage: humio login
  This command initializes your ~/.humioconfig file with
	the information needed to interact with a Humio cluster.
`
	return strings.TrimSpace(helpText)
}

func (f *LoginCommand) Synopsis() string {
	return "Setup your system for using Humio."
}

func (f *LoginCommand) Name() string { return "login" }

func (f *LoginCommand) Run(args []string) int {
	var addr, token string
	var err error

	f.Ui.Output("")
	f.Ui.Output("This will guide you through setting up the Humio CLI.")
	f.Ui.Output("")

	for true {
		f.Ui.Info("1. Which Humio instance should we talk to?")
		f.Ui.Output("")
		f.Ui.Output("If you are not using Humio Cloud enter the address out your Humio installation,")
		f.Ui.Output("e.g. http://localhost:8080/ or https://humio.example.com/")
		f.Ui.Output("")
		f.Ui.Output("Default: https://cloud.humio.com/ [Hit Enter]")
		addr, err = f.Ui.Ask("Humio Address:")

		if addr == "" {
			addr = "https://cloud.humio.com/"
		}

		if err != nil {
			f.Ui.Error(fmt.Sprintf("error readinf address: %s", err))
			return 1
		}

		// Make sure it is a valid URL and that
		// we always end in a slash.
		addrUrl, err := url.Parse(addr)

		if err != nil {
			f.Ui.Error("The valus must be a valid URL.")
		}

		addr = fmt.Sprintf("%v", addrUrl)

		clientConfig := api.DefaultConfig()
		clientConfig.Address = addr
		client, err := api.NewClient(clientConfig)

		if err != nil {
			f.Ui.Error(fmt.Sprintf("error initializing the http client: %s", err))
			return 1
		}

		status, err := client.Status()

		if err != nil {
			f.Ui.Error(fmt.Sprintf("Could not connect to the Humio server: %s", err))
			f.Ui.Error("Is the address connect and reachable?")
			continue
		}

		if status.Status != "ok" {
			f.Ui.Error(fmt.Sprintf("The server reported that is is malfunctioning, status: %s", status.Status))
			return 1
		}

		f.Ui.Output("")
		f.Ui.Info("==> Connection Successful")
		f.Ui.Output("")
		break
	}

	f.Ui.Output("")
	f.Ui.Info("2. Paste in your API Token")
	f.Ui.Output("")
	f.Ui.Output("To use Humio's CLI you will need to get a copy of your API Token.")
	f.Ui.Output("The API token can be found in your 'Account Settings' section of the UI.")
	f.Ui.Output("If you are running Humio without authorization just leave the API Token field empty.")
	f.Ui.Output("")
	f.Ui.Ask("We will now open the account page in a browser window. [Hit Any Key]")

	open.Start(fmt.Sprintf("%ssettings", addr))

	f.Ui.Output("")
	f.Ui.Output(fmt.Sprintf("If the browser did not open, you can manually visit:"))
	f.Ui.Output(fmt.Sprintf("%ssettings", addr))
	f.Ui.Output("")

	for true {
		token, err = f.Ui.AskSecret("API Token:")

		if err != nil {
			f.Ui.Error(fmt.Sprintf("error reading token: %s", err))
			return 1
		}

		// Create a new API client with the token
		config := api.DefaultConfig()
		config.Address = addr
		config.Token = token
		client, clientErr := api.NewClient(config)

		if clientErr != nil {
			f.Ui.Error(fmt.Sprintf("error initializing the http client: %s", clientErr))
			return 1
		}

		username, apiErr := client.Viewer().Username()

		if apiErr != nil {
			f.Ui.Error("authorization failed, try another token")
			continue
		}

		f.Ui.Output("")
		f.Ui.Info(fmt.Sprintf("==> Login successful '%s' 🎉", username))
		f.Ui.Output("")
		break
	}

	configFile, err := homedir.Expand("~/.humioconfig")

	if err != nil {
		f.Ui.Error(fmt.Sprintf("error writing settings file. Could not find home directory: %s", err))
		return 1
	}

	settingsTemplate := `
HUMIO_API_TOKEN="%s"
HUMIO_ADDR="%s"
`
	settings := fmt.Sprintf(settingsTemplate, token, addr)
	settings = strings.TrimSpace(settings)

	ioutil.WriteFile(configFile, []byte(settings), 0644)

	f.Ui.Info("==> Writing settings to: ~/.humioconfig")
	f.Ui.Output("")
	f.Ui.Output("Bye bye now!")
	f.Ui.Output("")

	return 0
}
