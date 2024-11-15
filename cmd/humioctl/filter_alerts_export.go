// Copyright © 2024 CrowdStrike
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

func newFilterAlertsExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <view> <filter-alert>",
		Short: "Export a filter alert <filter-alert> in <view> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			filterAlertName := args[1]
			client := NewApiClient(cmd)

			if outputName == "" {
				outputName = filterAlertName
			}

			filterAlerts, err := client.FilterAlerts().List(view)
			exitOnError(cmd, err, "Could not list filter alerts")

			var filterAlert api.FilterAlert
			for _, fa := range filterAlerts {
				if fa.Name == filterAlertName {
					filterAlert = fa
				}
			}

			if filterAlert.ID == "" {
				exitOnError(cmd, api.FilterAlertNotFound(filterAlertName), "Could not find filter alert")
			}

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

func newFilterAlertsExportAllCmd() *cobra.Command {
	var outputDirectory string

	cmd := cobra.Command{
		Use:   "export-all <view>",
		Short: "Export all filter alerts",
		Long:  `Export all filter alerts to yaml files with naming <sanitized-filter-alert-name>.yaml. All non-alphanumeric characters will be replaced with underscore.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			client := NewApiClient(cmd)

			var filterAlerts []api.FilterAlert
			filterAlerts, err := client.FilterAlerts().List(view)
			exitOnError(cmd, err, "Error fetching filter alerts")

			for _, filterAlert := range filterAlerts {
				yamlData, err := yaml.Marshal(&filterAlert)
				exitOnError(cmd, err, "Failed to serialize the filter alert")
				filterAlertFilename := sanitizeTriggerName(filterAlert.Name) + ".yaml"

				var outFilePath string
				if outputDirectory != "" {
					outFilePath = outputDirectory + "/" + filterAlertFilename
				} else {
					outFilePath = filterAlertFilename
				}

				err = os.WriteFile(outFilePath, yamlData, 0600)
				exitOnError(cmd, err, "Error saving the filter alert to file")
			}
		},
	}

	cmd.Flags().StringVarP(&outputDirectory, "outputDirectory", "d", "", "The file path where the filter alerts should be written. Defaults to current directory.")

	return &cmd
}
