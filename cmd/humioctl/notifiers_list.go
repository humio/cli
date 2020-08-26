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

	"github.com/spf13/cobra"
)

func newNotifiersListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list [flags] <view>",
		Short: "List all notifiers in a view.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			view := args[0]

			// Get the HTTP client
			client := NewApiClient(cmd)
			notifiers, err := client.Notifiers().List(view)

			if err != nil {
				return fmt.Errorf("Error fetching notifiers: %s", err)
			}

			var output []string
			output = append(output, "Name | Type")
			for i := 0; i < len(notifiers); i++ {
				notifier := notifiers[i]
				output = append(output, fmt.Sprintf("%v | %v", notifier.Name, notifier.Entity))
			}

			printTable(cmd, output)

			return nil
		},
	}

	return &cmd
}
