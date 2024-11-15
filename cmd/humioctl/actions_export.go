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

func newActionsExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <repo-or-view> <action>",
		Short: "Export an action <action> in <repo-or-view> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			actionName := args[1]
			client := NewApiClient(cmd)

			if outputName == "" {
				outputName = actionName
			}

			action, err := client.Actions().Get(repoOrViewName, actionName)
			exitOnError(cmd, err, "Error fetching action")

			yamlData, err := yaml.Marshal(&action)
			exitOnError(cmd, err, "Failed to serialize the action")

			outFilePath := outputName + ".yaml"
			err = os.WriteFile(outFilePath, yamlData, 0600)
			exitOnError(cmd, err, "Error saving the action file")
		},
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the action should be written. Defaults to ./<action-name>.yaml")

	return &cmd
}

func newActionsExportAllCmd() *cobra.Command {
	var outputDirectory string

	cmd := cobra.Command{
		Use:   "export-all <view>",
		Short: "Export all actions",
		Long:  `Export all actions to yaml files with naming <sanitized-action-name>.yaml. All non-alphanumeric characters will be replaced with underscore.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			client := NewApiClient(cmd)

			var actions []api.Action
			actions, err := client.Actions().List(view)
			exitOnError(cmd, err, "Error fetching actions")

			for _, action := range actions {
				yamlData, err := yaml.Marshal(&action)
				exitOnError(cmd, err, "Failed to serialize the action")
				actionFilename := sanitizeTriggerName(action.Name) + ".yaml"

				var outFilePath string
				if outputDirectory != "" {
					outFilePath = outputDirectory + "/" + actionFilename
				} else {
					outFilePath = actionFilename
				}

				err = os.WriteFile(outFilePath, yamlData, 0600)
				exitOnError(cmd, err, "Error saving the action to file")
			}
		},
	}

	cmd.Flags().StringVarP(&outputDirectory, "outputDirectory", "d", "", "The file path where the actions should be written. Defaults to current directory.")

	return &cmd
}
