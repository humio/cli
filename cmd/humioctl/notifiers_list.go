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
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newNotifiersListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list <repo-or-view>",
		Short: "List all notifiers in a repository or view.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			client := NewApiClient(cmd)

			notifiers, err := client.Notifiers().List(repoOrViewName)
			exitOnError(cmd, err, "Error fetching notifiers")

			var rows [][]format.Value
			for i := 0; i < len(notifiers); i++ {
				notifier := notifiers[i]
				rows = append(rows, []format.Value{format.String(notifier.Name), format.String(notifier.Entity)})
			}

			printOverviewTable(cmd, []string{"Name", "Type"}, rows)
		},
	}

	return &cmd
}
