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

func newFilterAlertsRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <view> <name>",
		Short: "Removes a filter alert.",
		Long:  `Removes the filter alert with name '<name>' in the view with name '<view>'.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]
			filterAlertName := args[1]
			client := NewApiClient(cmd)

			filterAlerts, err := client.FilterAlerts().List(viewName)
			exitOnError(cmd, err, "Could not list filter alerts")

			var filterAlert api.FilterAlert
			for _, fa := range filterAlerts {
				if fa.Name == filterAlertName {
					filterAlert = fa
				}
			}

			if filterAlert.ID == "" {
				exitOnError(cmd, api.FilterAlertNotFound(filterAlertName), "Could not find filter alert")
			}

			err = client.FilterAlerts().Delete(viewName, filterAlert.ID)
			exitOnError(cmd, err, "Could not remove alert")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed filter alert %q from view %q\n", filterAlertName, viewName)
		},
	}

	return cmd
}
