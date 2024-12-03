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

package main

import (
	"github.com/humio/cli/internal/format"
	"github.com/spf13/cobra"
)

func newParsersListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list <repo>",
		Short: "List all installed parsers in a repository.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			client := NewApiClient(cmd)

			parsers, err := client.Parsers().List(repo)
			exitOnError(cmd, err, "Error fetching parsers")

			var rows [][]format.Value
			for i := 0; i < len(parsers); i++ {
				parser := parsers[i]
				rows = append(rows, []format.Value{
					format.String(parser.Name),
					checkmark(!parser.IsBuiltIn),
					format.String(parser.ID),
				})
			}

			printOverviewTable(cmd, []string{"Name", "Custom", "ID"}, rows)
		},
	}

	return &cmd
}
