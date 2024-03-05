// Copyright © 2020 Humio Ltd.
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
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newFilterAlertsExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <view> <filter-alert>",
		Short: "Export an filter alert <filter-alert> in <view> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			filterAlertName := args[1]
			client := NewApiClient(cmd)

			if outputName == "" {
				outputName = filterAlertName
			}

			filterAlert, err := client.FilterAlerts().Get(view, filterAlertName)
			exitOnError(cmd, err, "Error fetching filter alert")

			yamlData, err := yaml.Marshal(&filterAlert)
			exitOnError(cmd, err, "Failed to serialize the filter alert")

			outFilePath := outputName + ".yaml"
			err = os.WriteFile(outFilePath, yamlData, 0600)
			exitOnError(cmd, err, "Error saving the filter alert file")
		},
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the filter alert should be written. Defaults to ./<filter-alert-name>.yaml")

	return &cmd
}
