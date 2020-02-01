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
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func newDashboardsInstallCmd() *cobra.Command {
	var content []byte
	var readErr error
	var sourceFile string

	cmd := cobra.Command{
		Use:   "install [flags] <repo>",
		Short: "Installs a dashboard from a source file to a repository",
		Long: `Install dashboard from a source file to a repsitory
				
					 This is in development, it WILL overwrite any existing dashboards
					 in the destination repository
					
					  $ humioctl dashboards install --source-file <source yaml file> <destination repo> 
`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check that we got the right number of arguments, all required

			if l := len(args); l != 1 {
				cmd.Println(fmt.Errorf("This command takes two arguments: <repo> --source-file <source yaml file>"))
				os.Exit(1)
			}

			if sourceFile == "" {
				cmd.Println(fmt.Errorf("This command takes two arguments: <repo> --source-file <source yaml file>"))
				os.Exit(1)
			}

			exitOnError(cmd, readErr, "Failed to load the dashboard")

			dashboard := api.Dashboard{}

			content, readErr = ioutil.ReadFile(sourceFile)
			yamlErr := yaml.Unmarshal(content, &dashboard)
			exitOnError(cmd, yamlErr, "The dashboard's format was invalid")

			dashboard.TemplateYaml = string(content)

			// Get the HTTP client
			client := NewApiClient(cmd)

			reposistoryName := args[0]

			installErr := client.Dashboards().Add(reposistoryName, &dashboard)
			exitOnError(cmd, installErr, "error installing dashboard")
			cmd.Println("Dashboard installed")
		},
	}

	cmd.Flags().StringVar(&sourceFile, "source-file", "", "The local file path to the source yaml.")

	return &cmd
}

// func getDashboardFromFile(filepath string) ([]byte, error) {
// 	return ioutil.readfile(filepath)
// }

// func getGithubParser(dashboardName string) ([]byte, error) {
// 	url := "https://raw.githubusercontent.com/humio/community/master/dashboards/" + dashboardName + ".yaml"
// 	return getUrlParser(url)
// }

// func getIUrlParser(url string) ([]byte, error) {
// 	response, err := http.Get(url)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer response.Body.Close()
// 	return ioutil.ReadAll(response.Body)
// }
