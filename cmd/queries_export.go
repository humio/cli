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
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"


	"github.com/spf13/cobra"
)

func newQueriesExportAllCmd() *cobra.Command {

	var exportPath string
	var outputPath string

	cmd := cobra.Command{
		Use:   "export-all [flags] <repo>",
		Short: "Export all queries form a <repo> to a files named the title of the query.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]

			if outputPath == "" {
				exportPath = "export/queries/"+repo+"/"
			}else{
				exportPath = outputPath + "/"
			}
			if _, err := os.Stat(exportPath); os.IsNotExist(err) {
				os.MkdirAll(exportPath, 0755)
			}
			// Get the HTTP client
			client := NewApiClient(cmd)

			queries, apiErr := client.Queries().GetAll(repo)
			if apiErr != nil {
				cmd.Println(fmt.Errorf("Error fetching queries for %s: %s", repo,apiErr))
				os.Exit(1)
			}
			// var objmap map[string]interface{}
			// var err error
			// var query []byte
			for _, query := range queries {
		  var name string
			name = strings.Replace(query.Name, " ", "_", -1)
		 	outFilePath := exportPath + name + ".json"

			var jsonOut []byte
			var jerror error
			jsonOut, jerror = json.Marshal(query)
			if jerror != nil {
				cmd.Println(fmt.Errorf("Error converting the query %s: %s", name, jerror))

			}
			writeErr := ioutil.WriteFile(outFilePath, jsonOut, 0644)
			if writeErr != nil {
				cmd.Println(fmt.Errorf("Error saving the query %s: %s", name, writeErr))
				os.Exit(1)
			}
				
			cmd.Println(fmt.Sprintf("Saved query %v to: %v", name, outFilePath))


			}

		},
	}

	cmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "The path where the queries should be written. Defaults to ./export/queries/<query_name>.json")

	return &cmd
}
