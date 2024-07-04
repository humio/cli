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
	"strings"

	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newFilterAlertsShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show <view> <name>",
		Short: "Show details about a filter alert in a view.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			name := args[1]
			client := NewApiClient(cmd)

			filterAlerts, err := client.FilterAlerts().List(view)
			exitOnError(cmd, err, "Could not list filter alert")

			var filterAlert api.FilterAlert
			for _, fa := range filterAlerts {
				if fa.Name == name {
					filterAlert = fa
				}
			}

			if filterAlert.ID == "" {
				exitOnError(cmd, api.FilterAlertNotFound(name), "Could not find filter alert")
			}

			details := [][]format.Value{
				{format.String("ID"), format.String(filterAlert.ID)},
				{format.String("Name"), format.String(filterAlert.Name)},
				{format.String("Enabled"), format.Bool(filterAlert.Enabled)},
				{format.String("Description"), format.String(filterAlert.Description)},
				{format.String("Query String"), format.String(filterAlert.QueryString)},
				{format.String("Labels"), format.String(strings.Join(filterAlert.Labels, ", "))},
				{format.String("Actions"), format.String(strings.Join(filterAlert.ActionNames, ", "))},
				{format.String("Throttle Time Seconds"), format.Int(filterAlert.ThrottleTimeSeconds)},
				{format.String("Throttle Field"), format.String(filterAlert.ThrottleField)},
				{format.String("Run As User ID"), format.String(filterAlert.RunAsUserID)},
				{format.String("Query Ownership Type"), format.String(filterAlert.QueryOwnershipType)},
			}

			printDetailsTable(cmd, details)
		},
	}

	return &cmd
}
