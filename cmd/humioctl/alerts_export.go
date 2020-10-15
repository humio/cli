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
	"fmt"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

func newAlertsExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <view> <alert>",
		Short: "Export an alert <alert> in <view> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			view := args[0]
			alertName := args[1]

			if outputName == "" {
				outputName = alertName
			}

			// Get the HTTP client
			client := NewApiClient(cmd)

			alert, apiErr := client.Alerts().Get(view, alertName)
			if apiErr != nil {
				return nil, fmt.Errorf("error fetching alert: %w", apiErr)
			}

			yamlData, yamlErr := yaml.Marshal(&alert)
			if yamlErr != nil {
				return nil, fmt.Errorf("failed to serialize the alert: %w", yamlErr)
			}
			outFilePath := outputName + ".yaml"

			writeErr := ioutil.WriteFile(outFilePath, yamlData, 0644)
			if writeErr != nil {
				return nil, fmt.Errorf("error saving the alert file: %w", writeErr)
			}

			return fmt.Sprintf("Alert %q exported to %s.", alertName, outFilePath), nil
		}),
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the alert should be written. Defaults to ./<alert-name>.yaml")

	return &cmd
}
