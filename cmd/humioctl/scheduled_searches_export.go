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
	"github.com/humio/cli/api"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newScheduledSearchesExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <view> <scheduled-search>",
		Short: "Export a scheduled search <scheduled-search> in <view> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			scheduledSearchName := args[1]
			client := NewApiClient(cmd)

			if outputName == "" {
				outputName = scheduledSearchName
			}

			scheduledSearches, err := client.ScheduledSearches().List(view)
			exitOnError(cmd, err, "Could not list scheduled searches")

			var scheduledSearch api.ScheduledSearch
			for _, ss := range scheduledSearches {
				if ss.Name == scheduledSearchName {
					scheduledSearch = ss
				}
			}

			if scheduledSearch.ID == "" {
				exitOnError(cmd, api.ScheduledSearchNotFound(scheduledSearchName), "Could not find scheduled search")
			}

			yamlData, err := yaml.Marshal(&scheduledSearch)
			exitOnError(cmd, err, "Failed to serialize the scheduled search")

			outFilePath := outputName + ".yaml"
			err = os.WriteFile(outFilePath, yamlData, 0600)
			exitOnError(cmd, err, "Error saving the scheduled search file")
		},
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the scheduled search should be written. Defaults to ./<scheduled-search-name>.yaml")

	return &cmd
}
