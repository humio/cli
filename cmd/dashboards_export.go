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
	"strings"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func newDashboardsExportAllCmd() *cobra.Command {

	var exportPath string
	var outputPath string

	cmd := cobra.Command{
		Use:   "export-all [flags] <repo>",
		Short: "Export all dashboards form a <repo> to a files named the title of the dashboard.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]

			if outputPath == "" {
				exportPath = "export/dashboards/"+repo+"/"
			}else{
				exportPath = outputPath + "/"
			}
			if _, err := os.Stat(exportPath); os.IsNotExist(err) {
				os.MkdirAll(exportPath, 0755)
			}
			// Get the HTTP client
			client := NewApiClient(cmd)

			dashboards, apiErr := client.Dashboards().GetAll(repo)
			if apiErr != nil {
				cmd.Println(fmt.Errorf("Error fetching dashboards for %s: %s", repo,apiErr))
				os.Exit(1)
			}
			// fmt.Println("%v", dashboards)
			for _, dashboard := range dashboards {
				var name string
				name = strings.Replace(dashboard.Name, " ", "_", -1)

				_, yamlErr := yaml.Marshal(dashboard.TemplateYaml)
				if yamlErr != nil {
					cmd.Println(fmt.Errorf("Failed to serialize the dashboard %s: %s", name, yamlErr))
					os.Exit(1)
				}
				outFilePath := exportPath + name + ".yaml"

				writeErr := ioutil.WriteFile(outFilePath, []byte(dashboard.TemplateYaml), 0644)
				if writeErr != nil {
					cmd.Println(fmt.Errorf("Error saving the dashboardfile %s: %s", name, writeErr))
					os.Exit(1)
				}
				
				cmd.Println(fmt.Sprintf("Saved dashboard %v to: %v", name, outFilePath))


			}

		},
	}

	cmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "The path where the dashboards should be written. Defaults to ./export/dashboards/<dashboard-name>.yaml")

	return &cmd
}
