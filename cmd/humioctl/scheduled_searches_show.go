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
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
	"strings"
)

func newScheduledSearchesShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show <view> <name>",
		Short: "Show details about a scheduled search in a view.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			name := args[1]
			client := NewApiClient(cmd)

			scheduledSearches, err := client.ScheduledSearches().List(view)
			exitOnError(cmd, err, "Could not list scheduled searches")

			var scheduledSearch api.ScheduledSearch
			for _, ss := range scheduledSearches {
				if ss.Name == name {
					scheduledSearch = ss
				}
			}

			if scheduledSearch.ID == "" {
				exitOnError(cmd, api.ScheduledSearchNotFound(name), "Could not find scheduled search")
			}

			details := [][]format.Value{
				{format.String("ID"), format.String(scheduledSearch.ID)},
				{format.String("Name"), format.String(scheduledSearch.Name)},
				{format.String("Description"), format.String(scheduledSearch.Description)},
				{format.String("Query String"), format.String(scheduledSearch.QueryString)},
				{format.String("Query Start"), format.String(scheduledSearch.QueryStart)},
				{format.String("Query End"), format.String(scheduledSearch.QueryEnd)},
				{format.String("Time Zone"), format.String(scheduledSearch.TimeZone)},
				{format.String("Schedule"), format.String(scheduledSearch.Schedule)},
				{format.String("Backfill Limit"), format.Int(scheduledSearch.BackfillLimit)},
				{format.String("Enabled"), format.Bool(scheduledSearch.Enabled)},
				{format.String("Actions"), format.String(strings.Join(scheduledSearch.ActionNames, ", "))},
				{format.String("Run As User ID"), format.String(scheduledSearch.RunAsUserID)},
				{format.String("Labels"), format.String(strings.Join(scheduledSearch.Labels, ", "))},
				{format.String("Query Ownership Type"), format.String(scheduledSearch.QueryOwnershipType)},
			}

			printDetailsTable(cmd, details)
		},
	}

	return &cmd
}
