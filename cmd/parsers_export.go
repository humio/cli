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

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func newParsersExportCmd() *cobra.Command {
	var outputName string

	cmd := cobra.Command{
		Use:   "export [flags] <repo> <parser>",
		Short: "Export a parser <parser> in <repo> to a file.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := args[0]
			parserName := args[1]

			if outputName == "" {
				outputName = parserName
			}

			// Get the HTTP client
			client := NewApiClient(cmd)

			parser, apiErr := client.Parsers().Get(repo, parserName)
			if apiErr != nil {
				return fmt.Errorf("Error fetching parsers: %s", apiErr)
			}

			yamlData, yamlErr := yaml.Marshal(&parser)
			if yamlErr != nil {
				return (fmt.Errorf("Failed to serialize the parser: %s", yamlErr))
			}
			outFilePath := outputName + ".yaml"

			writeErr := ioutil.WriteFile(outFilePath, yamlData, 0644)
			if writeErr != nil {
				return (fmt.Errorf("Error saving the parser file: %s", writeErr))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputName, "output", "o", "", "The file path where the parser should be written. Defaults to ./<parser-name>.yaml")

	return &cmd
}
