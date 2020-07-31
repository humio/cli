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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newAlertsListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list [flags] <view>",
		Short: "List all alerts in a view.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			view := args[0]

			// Get the HTTP client
			client := NewApiClient(cmd)
			alerts, err := client.Alerts().List(view)

			if err != nil {
				return fmt.Errorf("Error fetching alerts: %s", err)
			}

			var output []string
			output = append(output, "Name | Enabled | Description | Notifiers")
			for i := 0; i < len(alerts); i++ {
				alert := alerts[i]
				var notifierNames []string
				for _, notifierID := range alert.Notifiers {
					notifier, err := client.Notifiers().GetByID(view, notifierID)
					if err != nil {
						return fmt.Errorf("could not get details for notifier with id %s: %v", notifierID, err)
					}
					notifierNames = append(notifierNames, notifier.Name)
				}
				output = append(output, fmt.Sprintf("%v | %v | %v | %v", alert.Name, !alert.Silenced, alert.Description, strings.Join(notifierNames, ", ")))
			}

			printTable(cmd, output)

			return nil
		},
	}

	return &cmd
}
