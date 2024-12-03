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

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newAlertsExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <view> <alert>",
		Short: "Export an alert <alert> in <view> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			alertName := args[1]
			client := NewApiClient(cmd)

			if outputName == "" {
				outputName = alertName
			}

			alert, err := client.Alerts().Get(view, alertName)
			exitOnError(cmd, err, "Error fetching alert")

			yamlData, err := yaml.Marshal(&alert)
			exitOnError(cmd, err, "Failed to serialize the alert")

			outFilePath := outputName + ".yaml"
			err = os.WriteFile(outFilePath, yamlData, 0600)
			exitOnError(cmd, err, "Error saving the alert file")
		},
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the alert should be written. Defaults to ./<alert-name>.yaml")

	return &cmd
}

func newAlertsExportAllCmd() *cobra.Command {
	var outputDirectory string

	cmd := cobra.Command{
		Use:   "export-all <view>",
		Short: "Export all alerts",
		Long:  `Export all alerts to yaml files with naming <sanitized-alert-name>.yaml. All non-alphanumeric characters will be replaced with underscore.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			client := NewApiClient(cmd)

			var alerts []api.Alert
			alerts, err := client.Alerts().List(view)
			exitOnError(cmd, err, "Error fetching alerts")

			for i := range alerts {
				yamlData, err := yaml.Marshal(&alerts[i])
				exitOnError(cmd, err, "Failed to serialize the alert")
				alertFilename := sanitizeTriggerName(alerts[i].Name) + ".yaml"

				var outFilePath string
				if outputDirectory != "" {
					outFilePath = outputDirectory + "/" + alertFilename
				} else {
					outFilePath = alertFilename
				}

				err = os.WriteFile(outFilePath, yamlData, 0600)
				exitOnError(cmd, err, "Error saving the alert to file")
			}
		},
	}

	cmd.Flags().StringVarP(&outputDirectory, "outputDirectory", "d", "", "The file path where the alerts should be written. Defaults to current directory.")

	return &cmd
}
