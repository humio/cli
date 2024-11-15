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

func newAggregateAlertsExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <view> <aggregate-alert>",
		Short: "Export an aggregate alert <aggregate-alert> in <view> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			aggregateAlertName := args[1]
			client := NewApiClient(cmd)

			if outputName == "" {
				outputName = aggregateAlertName
			}

			aggregateAlerts, err := client.AggregateAlerts().List(view)
			exitOnError(cmd, err, "Could not list aggregate alerts")

			var aggregateAlert api.AggregateAlert
			for _, fa := range aggregateAlerts {
				if fa.Name == aggregateAlertName {
					aggregateAlert = fa
				}
			}

			if aggregateAlert.ID == "" {
				exitOnError(cmd, api.AggregateAlertNotFound(aggregateAlertName), "Could not find aggregate alert")
			}

			yamlData, err := yaml.Marshal(&aggregateAlert)
			exitOnError(cmd, err, "Failed to serialize the aggregate alert")

			outFilePath := outputName + ".yaml"
			err = os.WriteFile(outFilePath, yamlData, 0600)
			exitOnError(cmd, err, "Error saving the aggregate alert file")
		},
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the aggregate alert should be written. Defaults to ./<aggregate-alert-name>.yaml")

	return &cmd
}

func newAggregateAlertsExportAllCmd() *cobra.Command {
	var outputDirectory string

	cmd := cobra.Command{
		Use:   "export-all <view>",
		Short: "Export all aggregate alerts",
		Long:  `Export all aggregate alerts to yaml files with naming <sanitized-aggregate-alert-name>.yaml. All non-alphanumeric characters will be replaced with underscore.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			client := NewApiClient(cmd)

			var aggregateAlerts []api.AggregateAlert
			aggregateAlerts, err := client.AggregateAlerts().List(view)
			exitOnError(cmd, err, "Error fetching aggregate alerts")

			for i := range aggregateAlerts {
				yamlData, err := yaml.Marshal(&aggregateAlerts[i])
				exitOnError(cmd, err, "Failed to serialize the aggregate alert")
				alertFilename := sanitizeTriggerName(aggregateAlerts[i].Name) + ".yaml"

				var outFilePath string
				if outputDirectory != "" {
					outFilePath = outputDirectory + "/" + alertFilename
				} else {
					outFilePath = alertFilename
				}

				err = os.WriteFile(outFilePath, yamlData, 0600)
				exitOnError(cmd, err, "Error saving the aggregate alert to file")
			}
		},
	}

	cmd.Flags().StringVarP(&outputDirectory, "outputDirectory", "d", "", "The file path where the aggregate alerts should be written. Defaults to current directory.")

	return &cmd
}
