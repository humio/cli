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
	"os"
	"path"
	"strconv"

	"github.com/humio/cli/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, tokenFile, token, address, caCertificateFile, profileFlag string
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
				fmt.Println(fmt.Sprintf("humioctl %s (%s on %s)", version, commit, date))
				os.Exit(0)
			}

			// If no token or address flags are passed
			// and no configuration file exists, run login.
			if viper.GetString("token") == "" && viper.GetString("address") == "" {
				if err := newWelcomeCmd().Execute(); err != nil {
					fmt.Println(fmt.Errorf("error printing welcome message: %v", err))
				}

			} else {
				if err := cmd.Help(); err != nil {
					fmt.Println(fmt.Errorf("error printing help: %v", err))
				}
			}
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetOutput(os.Stdout)
		},
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&profileFlag, "profile", "u", "", "Name of the config profile to use")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.humio/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "The API token to user when talking to Humio. Overrides the value in your config file.")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token-file", "", "File path to a file containing the API token. Overrides the value in your config file and the value of --token.")
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "The HTTP address of the Humio cluster. Overrides the value in your config file.")
	rootCmd.PersistentFlags().StringVar(&caCertificateFile, "ca-certificate-file", "", "File path to a file containing the CA certificate in PEM format. Overrides the value in your config file.")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "By default, all encrypted connections will verify that the hostname in the TLS certificate matches the name from the URL. Set this to true to ignore hostname validation.")

	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("token-file", rootCmd.PersistentFlags().Lookup("token-file"))
	viper.BindPFlag("ca-certificate-file", rootCmd.PersistentFlags().Lookup("ca-certificate-file"))
	viper.BindPFlag("insecure", rootCmd.PersistentFlags().Lookup("insecure"))

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
	viper.ReadInConfig()

	// If the user has specified a profile flag, load it.
	if profileFlag != "" {
		profile, loadErr := loadProfile(profileFlag)
		if loadErr != nil {
			fmt.Println(fmt.Errorf("failed to load profile: %s", loadErr))
			os.Exit(1)
		}

		// Explicitly bound address or token have precedence
		if address == "" {
			viper.Set("address", profile.address)
		}
		if token == "" {
			viper.Set("token", profile.token)
		}
		if caCertificateFile == "" {
			viper.Set("ca-certificate", profile.caCertificate)
		}
		if insecure {
			viper.Set("insecure", strconv.FormatBool(insecure))
		}
	}

	if tokenFile != "" {
		tokenFileContent, tokenFileErr := ioutil.ReadFile(tokenFile)
		if tokenFileErr != nil {
			fmt.Println(fmt.Sprintf("error loading token file: %s", tokenFileErr))
			os.Exit(1)
		}
		viper.Set("token", string(tokenFileContent))
	}

	if caCertificateFile != "" {
		caCertificateFileContent, caCertificateFileErr := ioutil.ReadFile(caCertificateFile)
		if caCertificateFileErr != nil {
			fmt.Println(fmt.Sprintf("error loading CA certificate file: %s", caCertificateFileErr))
			os.Exit(1)
		}
		viper.Set("ca-certificate", string(caCertificateFileContent))
	}

	if insecure {
		viper.Set("insecure", insecure)
	}
}

func NewApiClient(cmd *cobra.Command) *api.Client {
	client, err := newApiClientE(cmd)

	if err != nil {
		fmt.Println(fmt.Errorf("Error creating HTTP client: %s", err))
		os.Exit(1)
	}

	return client
}

func newApiClientE(cmd *cobra.Command) (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = viper.GetString("address")
	config.Token = viper.GetString("token")
	config.CACertificate = []byte(viper.GetString("ca-certificate"))
	config.Insecure = viper.GetBool("insecure")

	return api.NewClient(config)
}

func main() {
	SetVersion(version, commit, date)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
