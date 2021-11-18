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

package humioctl

import (
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

func newNotifiersExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <repo-or-view> <notifier>",
		Short: "Export a notifier <notifier> in <repo-or-view> to a file.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			notifierName := args[1]
			client := NewApiClient(cmd)

			if outputName == "" {
				outputName = notifierName
			}

			notifier, err := client.Notifiers().Get(repoOrViewName, notifierName)
			helpers.ExitOnError(cmd, err, "Error fetching notifier")

			yamlData, err := yaml.Marshal(&notifier)
			helpers.ExitOnError(cmd, err, "Failed to serialize the notifier")

			outFilePath := outputName + ".yaml"
			err = ioutil.WriteFile(outFilePath, yamlData, 0600)
			helpers.ExitOnError(cmd, err, "Error saving the notifier file")
		},
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the notifier should be written. Defaults to ./<notifier-name>.yaml")

	return &cmd
}
