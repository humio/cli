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
	"fmt"

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
)

func newAggregateAlertsRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <view> <name>",
		Short: "Removes an aggregate alert.",
		Long:  `Removes the aggregate alert with name '<name>' in the view with name '<view>'.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]
			aggregateAlertName := args[1]
			client := NewApiClient(cmd)

			aggregateAlerts, err := client.AggregateAlerts().List(viewName)
			exitOnError(cmd, err, "Could not list aggregate alerts")

			var aggregateAlert api.AggregateAlert
			for _, fa := range aggregateAlerts {
				if fa.Name == aggregateAlertName {
					aggregateAlert = fa
				}
			}

			if aggregateAlert.ID == "" {
				exitOnError(cmd, api.AggregateAlertNotFound(aggregateAlertName), "Could not find aggregate alert")
			}

			err = client.AggregateAlerts().Delete(viewName, aggregateAlert.ID)
			exitOnError(cmd, err, "Could not remove alert")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed aggregate alert %q from view %q\n", aggregateAlertName, viewName)
		},
	}

	return cmd
}
