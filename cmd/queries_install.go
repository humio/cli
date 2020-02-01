// // Copyright © 2018 Humio Ltd.
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //     http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"encoding/json"
)

func newQueryInstallCmd() *cobra.Command {
	var content []byte
	var readErr error
	var sourceFile string

	cmd := cobra.Command{
		Use:   "install [flags] <repo> <dashboard yaml file>",
		Short: "Installs a query from a source file to a repository",
		Long: `Install  a query from a source file to a repository
				
					
					  $ humioctl query install --source-file <source json file> <destination repo>
`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check that we got the right number of arguments, all required

			if l := len(args); l != 1 {
				cmd.Println(fmt.Errorf("This command takes two arguments: --source-file <source yaml file> <repo>"))
				os.Exit(1)
			}

			if sourceFile == "" {
				cmd.Println(fmt.Errorf("This command takes two arguments: --source-file <source yaml file> <repo>"))
				os.Exit(1)
			}

			exitOnError(cmd, readErr, "Failed to load the query")

			query := api.QueryJSON{}

			jsonFile, err := os.Open(sourceFile)
			if err != nil {
				exitOnError(cmd, readErr, "Failed to load the query")
			}
			content, readErr = ioutil.ReadAll(jsonFile)

			jsonErr := json.Unmarshal(content, &query)
			exitOnError(cmd, jsonErr, "The query's format was invalid")


			// Get the HTTP client
			client := NewApiClient(cmd)

			reposistoryName := args[0]

			installErr := client.Queries().Add(reposistoryName, &query)
			exitOnError(cmd, installErr, "error installing query")
			cmd.Println("Query installed")
		},
	}

	cmd.Flags().StringVar(&sourceFile, "source-file", "", "The local file path to the source yaml.")

	return &cmd
}

// // func getDashboardFromFile(filepath string) ([]byte, error) {
// // 	return ioutil.readfile(filepath)
// // }

// // func getGithubParser(dashboardName string) ([]byte, error) {
// // 	url := "https://raw.githubusercontent.com/humio/community/master/dashboards/" + dashboardName + ".yaml"
// // 	return getUrlParser(url)
// // }

// // func getIUrlParser(url string) ([]byte, error) {
// // 	response, err := http.Get(url)

// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	defer response.Body.Close()
// // 	return ioutil.ReadAll(response.Body)
// // }
