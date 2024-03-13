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
	"strings"

	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newFilterAlertsListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list <view>",
		Short: "List all filter alerts in a view.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			client := NewApiClient(cmd)

			filterAlerts, err := client.FilterAlerts().List(view)
			exitOnError(cmd, err, "Error fetching filter alerts")

			var rows = make([][]format.Value, len(filterAlerts))
			for i := range filterAlerts {
				filterAlert := filterAlerts[i]

				var actionNames = make([]string, len(filterAlert.ActionNames))
				for j := range filterAlert.ActionNames {
					actionNames[j] = filterAlert.ActionNames[j]
				}

				rows[i] = []format.Value{
					format.String(filterAlert.ID),
					format.String(filterAlert.Name),
					format.Bool(filterAlert.Enabled),
					format.String(filterAlert.Description),
					format.String(strings.Join(actionNames, ", ")),
					format.String(strings.Join(filterAlert.Labels, ", ")),
					format.String(filterAlert.RunAsUserID),
					format.String(filterAlert.QueryOwnershipType),
				}
			}

			printOverviewTable(cmd, []string{"ID", "Name", "Enabled", "Description", "Actions", "Labels", "Run As User ID", "Query Ownership Type"}, rows)
		},
	}

	return &cmd
}
