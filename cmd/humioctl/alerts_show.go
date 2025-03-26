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
	"strings"

	"github.com/humio/cli/internal/format"
	"github.com/spf13/cobra"
)

func newAlertsShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show <repo-or-view> <name>",
		Short: "Show details about an alert in a repository or view.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			name := args[1]
			client := NewApiClient(cmd)

			alert, err := client.Alerts().Get(repoOrViewName, name)
			exitOnError(cmd, err, "Error fetching alert")

			details := [][]format.Value{
				{format.String("ID"), format.String(alert.ID)},
				{format.String("Name"), format.String(alert.Name)},
				{format.String("Enabled"), format.Bool(alert.Enabled)},
				{format.String("Description"), format.StringPtr(alert.Description)},
				{format.String("Query Start"), format.String(alert.QueryStart)},
				{format.String("Query String"), format.String(alert.QueryString)},
				{format.String("Labels"), format.String(strings.Join(alert.Labels, ", "))},
				{format.String("Throttle Time Millis"), format.Int(alert.ThrottleTimeMillis)},
				{format.String("Is Starred"), format.Bool(alert.IsStarred)},
				{format.String("Last Error"), format.StringPtr(alert.LastError)},
				{format.String("Throttle Field"), format.StringPtr(alert.ThrottleField)},
				{format.String("Time Of Last Trigger"), format.IntPtr(alert.TimeOfLastTrigger)},
				{format.String("Run As User ID"), format.String(alert.RunAsUserID)},
				{format.String("Query Ownership Type"), format.String(alert.QueryOwnershipType)},
			}

			printDetailsTable(cmd, details)
		},
	}

	return &cmd
}
