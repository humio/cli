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
	"strings"
)

func newAlertsListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list [flags] <view>",
		Short: "List all alerts in a view.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			client := NewApiClient(cmd)

			alerts, err := client.Alerts().List(view)
			exitOnError(cmd, err, "Error fetching alerts")

			notifiers, err := client.Notifiers().List(view)
			exitOnError(cmd, err, "Unable to fetch notifier details")

			var notifierMap = map[string]string{}
			for _, notifier := range notifiers {
				notifierMap[notifier.ID] = notifier.Name
			}

			var rows [][]format.Value
			for i := 0; i < len(alerts); i++ {
				alert := alerts[i]
				var notifierNames []string
				for _, notifierID := range alert.Notifiers {
					notifierNames = append(notifierNames, notifierMap[notifierID])
				}
				rows = append(rows, []format.Value{
					format.String(alert.Name),
					format.Bool(!alert.Silenced),
					format.String(alert.Description),
					format.String(strings.Join(notifierNames, ", "))})
			}

			printOverviewTable(cmd, []string{"Name", "Enabled", "Description", "Notifiers"}, rows)
		},
	}

	return &cmd
}
