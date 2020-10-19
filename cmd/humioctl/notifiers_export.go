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

func newNotifiersExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <view> <notifier>",
		Short: "Export a notifier <notifier> in <view> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			view := args[0]
			notifierName := args[1]

			if outputName == "" {
				outputName = notifierName
			}

			// Get the HTTP client
			client := NewApiClient(cmd)

			notifier, apiErr := client.Notifiers().Get(view, notifierName)
			if apiErr != nil {
				return nil, fmt.Errorf("error fetching notifier: %w", apiErr)
			}

			yamlData, yamlErr := yaml.Marshal(&notifier)
			if yamlErr != nil {
				return nil, fmt.Errorf("failed to serialize the notifier: %w", yamlErr)
			}
			outFilePath := outputName + ".yaml"

			writeErr := ioutil.WriteFile(outFilePath, yamlData, 0644)
			if writeErr != nil {
				return nil, fmt.Errorf("Error saving the notifier file: %w", writeErr)
			}

			return nil, nil
		}),
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the notifier should be written. Defaults to ./<notifier-name>.yaml")

	return &cmd
}
