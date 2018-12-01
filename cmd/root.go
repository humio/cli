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

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/humio/cli/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, tokenFile, token, address string

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd = &cobra.Command{
		Use:   "humio [subcommand] [flags] [arguments]",
		Short: "A management CLI for Humio.",
		Long: `
Environment Setup:

  $ humio login

Sending Data:

  Humio's CLI is not a replacement for fully-featured data-shippers like
  LogStash, FileBeat or MetricBeat. It can be handy to easily send logs
  to Humio, e.g examine a local log file or test a parser on test input.

To stream the content of "/var/log/system.log" data to Humio:

  $ tail -f /var/log/system.log | humio ingest -o

or

  $ humio ingest -o --tail=/var/log/system.log

Common Management Commands:
  users <subcommand>
  parsers <subcommand>
  views <subcommand>
		`,
		Run: func(cmd *cobra.Command, args []string) {
			// If no token or address flags are passed
			// and no configuration file exists, run login.
			if viper.GetString("token") == "" && viper.GetString("address") == "" {
				newLoginCmd().Execute()
			} else {
				err := cmd.Help()
				if err != nil {
					fmt.Println(fmt.Errorf("error printing help: %s", err))
				}
			}
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Name() != "login" {
				cmd.Println()
				cmd.Println("Humio Address:", viper.GetString("address"))
				cmd.Println()
			}
		},
	}

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.humio/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "The API token to user when talking to Humio. Overrides the value in your config file.")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token-file", "", "File path to a file containing the API token. Overrides the value in your config file and the value of --token.")
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "The HTTP address of the Humio cluster. Overrides the value in your config file.")

	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("token-file", rootCmd.PersistentFlags().Lookup("token-file"))

	rootCmd.AddCommand(newUsersCmd())
	rootCmd.AddCommand(newParsersCmd())
	rootCmd.AddCommand(newIngestCmd())
	rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(newIngestTokensCmd())
	rootCmd.AddCommand(newViewsCmd())
	rootCmd.AddCommand(newCompletionCmd())
	rootCmd.AddCommand(newLicenseCmd())
	rootCmd.AddCommand(newReposCmd())
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

	if tokenFile != "" {
		tokenFileContent, tokenFileErr := ioutil.ReadFile(tokenFile)
		if tokenFileErr != nil {
			fmt.Println(fmt.Sprintf("error loading token file: %s", tokenFileErr))
			os.Exit(1)
		}
		viper.Set("token", string(tokenFileContent))
	}
}

func NewApiClient(cmd *cobra.Command) *api.Client {
	client, err := NewApiClientE(cmd)

	if err != nil {
		fmt.Println(fmt.Errorf("Error creating HTTP client: %s", err))
		os.Exit(1)
	}

	return client
}

func NewApiClientE(cmd *cobra.Command) (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = viper.GetString("address")
	config.Token = viper.GetString("token")

	return api.NewClient(config)
}
