// Copyright © 2025 CrowdStrike
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
	"strconv"
	"strings"

	"github.com/humio/cli/internal/format"
	"github.com/spf13/cobra"
)

func newScheduledSearchesV2ListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list <view>",
		Short: "List all scheduled searches in a view.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			client := NewApiClient(cmd)

			scheduledSearches, err := client.ScheduledSearchesV2().List(view)
			exitOnError(cmd, err, "Error fetching scheduled searches")

			var rows = make([][]format.Value, len(scheduledSearches))
			for i := range scheduledSearches {
				scheduledSearch := scheduledSearches[i]

				rows[i] = []format.Value{
					format.String(scheduledSearch.ID),
					format.String(scheduledSearch.Name),
					format.StringPtr(scheduledSearch.Description),
					format.String(strconv.FormatInt(scheduledSearch.SearchIntervalSeconds, 10)),
					format.Int64Ptr(scheduledSearch.SearchIntervalOffsetSeconds),
					format.String(scheduledSearch.QueryTimestampType),
					format.Int64Ptr(scheduledSearch.MaxWaitTimeSeconds),
					format.IntPtr(scheduledSearch.BackfillLimitV2),
					format.String(scheduledSearch.TimeZone),
					format.String(scheduledSearch.Schedule),
					format.String(strings.Join(scheduledSearch.ActionNames, ", ")),
					format.String(strings.Join(scheduledSearch.Labels, ", ")),
					format.Bool(scheduledSearch.Enabled),
					format.String(scheduledSearch.OwnershipRunAsID),
					format.String(scheduledSearch.QueryOwnershipType),
				}
			}

			printOverviewTable(cmd, []string{
				"ID",
				"Name",
				"Description",
				"Search Interval Seconds",
				"Search Interval Offset Seconds",
				"Query Timestamp Type",
				"Max Wait Time",
				"Backfill Limit",
				"Time Zone",
				"Schedule",
				"Action Names",
				"Labels",
				"Enabled",
				"Run As User ID",
				"Query Ownership Type",
			}, rows)
		},
	}

	return &cmd
}
