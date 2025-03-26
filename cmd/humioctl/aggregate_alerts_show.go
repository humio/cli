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

	"github.com/humio/cli/internal/api"
	"github.com/humio/cli/internal/format"
	"github.com/spf13/cobra"
)

func newAggregateAlertsShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show <view> <name>",
		Short: "Show details about an aggregate alert in a view.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			name := args[1]
			client := NewApiClient(cmd)

			aggregateAlerts, err := client.AggregateAlerts().List(view)
			exitOnError(cmd, err, "Could not list aggregate alert")

			var aggregateAlert api.AggregateAlert
			for _, fa := range aggregateAlerts {
				if fa.Name == name {
					aggregateAlert = fa
				}
			}

			if aggregateAlert.ID == "" {
				exitOnError(cmd, api.AggregateAlertNotFound(name), "Could not find aggregate alert")
			}

			details := [][]format.Value{
				{format.String("ID"), format.String(aggregateAlert.ID)},
				{format.String("Name"), format.String(aggregateAlert.Name)},
				{format.String("Description"), format.StringPtr(aggregateAlert.Description)},
				{format.String("Query String"), format.String(aggregateAlert.QueryString)},
				{format.String("Search Interval Seconds"), format.Int64(aggregateAlert.SearchIntervalSeconds)},
				{format.String("Actions"), format.String(strings.Join(aggregateAlert.ActionNames, ", "))},
				{format.String("Labels"), format.String(strings.Join(aggregateAlert.Labels, ", "))},
				{format.String("Enabled"), format.Bool(aggregateAlert.Enabled)},
				{format.String("Throttle Field"), format.StringPtr(aggregateAlert.ThrottleField)},
				{format.String("Throttle Time Seconds"), format.Int64(aggregateAlert.ThrottleTimeSeconds)},
				{format.String("Query Timestamp Type"), format.String(aggregateAlert.QueryTimestampType)},
				{format.String("Trigger Mode"), format.String(aggregateAlert.TriggerMode)},
				{format.String("Run As User ID"), format.String(aggregateAlert.OwnershipRunAsID)},
				{format.String("Query Ownership Type"), format.String(aggregateAlert.QueryOwnershipType)},
			}

			printDetailsTable(cmd, details)
		},
	}

	return &cmd
}
