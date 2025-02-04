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
	"fmt"

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
)

func newScheduledSearchesV2RemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <view> <name>",
		Short: "Removes a scheduled search.",
		Long:  `Removes the scheduled search with name '<name>' in the view with name '<view>'.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]
			scheduledSearchName := args[1]
			client := NewApiClient(cmd)

			scheduledSearches, err := client.ScheduledSearchesV2().List(viewName)
			exitOnError(cmd, err, "Could not list scheduled searches")

			var scheduledSearch api.ScheduledSearchV2
			for _, ss := range scheduledSearches {
				if ss.Name == scheduledSearchName {
					scheduledSearch = ss
				}
			}

			if scheduledSearch.ID == "" {
				exitOnError(cmd, api.ScheduledSearchNotFound(scheduledSearchName), "Could not find scheduled search")
			}

			err = client.ScheduledSearchesV2().Delete(viewName, scheduledSearch.ID)
			exitOnError(cmd, err, "Could not remove scheduled search")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed scheduled search %q from view %q\n", scheduledSearchName, viewName)
		},
	}

	return cmd
}
