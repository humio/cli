// Copyright © 2018 Humio Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/humio/cli/cmd/humioctl/internal/viperkey"

	"github.com/humio/cli/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, tokenFile, token, address, caCertificateFile, profileFlag, proxyOrganization string
var insecure bool

var printVersion bool

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use:   "humioctl [subcommand] [flags] [arguments]",
		Short: "A management CLI for Humio.",
		Long: `
Sending Data:

  Humio's CLI is not a replacement for fully-featured data-shippers like
  LogStash, FileBeat or MetricBeat. It can be handy to easily send logs
  to Humio, e.g examine a local log file or test a parser on test input.

To stream the content of "/var/log/system.log" data to Humio:

  $ tail -f /var/log/system.log | humioctl ingest -o

or

  $ humioctl ingest -o --tail=/var/log/system.log

Common Management Commands:
  users <subcommand>
  parsers <subcommand>
  views <subcommand>
	status
		`,
		Run: func(cmd *cobra.Command, args []string) {

			if printVersion {
				cmd.Printf("humioctl %s (%s on %s)\n", version, commit, date)
				os.Exit(0)
			}

			// If no token or address flags are passed
			// and no configuration file exists, run login.
			if viper.GetString(viperkey.Token) == "" && viper.GetString(viperkey.Address) == "" {
				if err := newWelcomeCmd().Execute(); err != nil {
					fmt.Println(fmt.Errorf("error printing welcome message: %v", err))
				}

			} else {
				if err := cmd.Help(); err != nil {
					fmt.Println(fmt.Errorf("error printing help: %v", err))
				}
			}
		},
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&profileFlag, "profile", "u", "", "Name of the config profile to use")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.humio/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "The API token to use when talking to Humio. Overrides the value in your config file.")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token-file", "", "File path to a file containing the API token. Overrides the value in your config file and the value of --token.")
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "The HTTP address of the Humio cluster. Overrides the value in your config file.")
	rootCmd.PersistentFlags().StringVar(&caCertificateFile, "ca-certificate-file", "", "File path to a file containing the CA certificate in PEM format. Overrides the value in your config file.")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "By default, all encrypted connections will verify that the hostname in the TLS certificate matches the name from the URL. Set this to true to ignore hostname validation.")
	rootCmd.PersistentFlags().StringVar(&proxyOrganization, "proxy-organization", "", "Commands are executed in the specified organization.")

	_ = viper.BindPFlag(viperkey.Address, rootCmd.PersistentFlags().Lookup("address"))
	_ = viper.BindPFlag(viperkey.Token, rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag(viperkey.TokenFile, rootCmd.PersistentFlags().Lookup("token-file"))
	_ = viper.BindPFlag(viperkey.CACertificateFile, rootCmd.PersistentFlags().Lookup("ca-certificate-file"))
	_ = viper.BindPFlag(viperkey.Insecure, rootCmd.PersistentFlags().Lookup("insecure"))
	_ = viper.BindPFlag(viperkey.ProxyOrganization, rootCmd.PersistentFlags().Lookup("proxy-organization"))

	rootCmd.Flags().BoolVarP(&printVersion, "version", "v", false, "Print the client version")

	rootCmd.AddCommand(newUsersCmd())
	rootCmd.AddCommand(newParsersCmd())
	rootCmd.AddCommand(newIngestCmd())
	rootCmd.AddCommand(newProfilesCmd())
	rootCmd.AddCommand(newIngestTokensCmd())
	rootCmd.AddCommand(newViewsCmd())
	rootCmd.AddCommand(newCompletionCmd())
	rootCmd.AddCommand(newLicenseCmd())
	rootCmd.AddCommand(newReposCmd())
	rootCmd.AddCommand(newSearchCmd())
	rootCmd.AddCommand(newStatusCmd())
	rootCmd.AddCommand(newHealthCmd())
	rootCmd.AddCommand(newClusterCmd())
	rootCmd.AddCommand(newNotifiersCmd())
	rootCmd.AddCommand(newAlertsCmd())
	rootCmd.AddCommand(newPackagesCmd())
	rootCmd.AddCommand(newGroupsCmd())
	rootCmd.AddCommand(newTransferCmd())
	rootCmd.AddCommand(newFilesCmd())

	// Hidden Commands
	rootCmd.AddCommand(newWelcomeCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cfgFile = path.Join(home, ".humio", "config.yaml")
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")
	}

	viper.SetEnvPrefix("HUMIO")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()

	// If the user has specified a profile flag, load it.
	if profileFlag != "" {
		profile, err := loadProfile(profileFlag)
		exitOnError(rootCmd, err, "Failed to load profile")

		// Explicitly bound address or token have precedence
		if address == "" {
			viper.Set(viperkey.Address, profile.address)
		}
		if token == "" {
			viper.Set(viperkey.Token, profile.token)
		}
		if caCertificateFile == "" {
			viper.Set(viperkey.CACertificate, profile.caCertificate)
		}
		if !insecure {
			viper.Set(viperkey.Insecure, profile.insecure)
		}
	}

	if tokenFile != "" {
		// #nosec G304
		tokenFileContent, err := ioutil.ReadFile(tokenFile)
		exitOnError(rootCmd, err, "Error loading token file")
		viper.Set(viperkey.Token, string(tokenFileContent))
	}

	if caCertificateFile != "" {
		// #nosec G304
		caCertificateFileContent, err := ioutil.ReadFile(caCertificateFile)
		exitOnError(rootCmd, err, "Error loading CA certificate file")
		viper.Set(viperkey.CACertificate, string(caCertificateFileContent))
	}

	if insecure {
		viper.Set(viperkey.Insecure, insecure)
	}
}

func NewApiClient(cmd *cobra.Command, opts ...func(config *api.Config)) *api.Client {
	client, err := newApiClientE(opts...)
	exitOnError(cmd, err, "Error creating HTTP client")
	return client
}

func newApiClientE(opts ...func(config *api.Config)) (*api.Client, error) {
	config := api.DefaultConfig()
	addr := viper.GetString(viperkey.Address)
	parsedURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	config.Address = parsedURL
	config.Token = viper.GetString(viperkey.Token)
	config.CACertificatePEM = viper.GetString(viperkey.CACertificate)
	config.Insecure = viper.GetBool(viperkey.Insecure)
	config.ProxyOrganization = viper.GetString(viperkey.ProxyOrganization)
	config.UserAgent = fmt.Sprintf("humioctl/%s (%s on %s)", version, commit, date)

	for _, opt := range opts {
		opt(&config)
	}

	return api.NewClient(config), nil
}

func main() {
	SetVersion(version, commit, date)
	err := rootCmd.Execute()
	exitOnError(rootCmd, err, "Unable to execute command")
}
