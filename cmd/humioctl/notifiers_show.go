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

func newNotifiersShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show [flags] <view> <name>",
		Short: "Show details about a notifier in a view.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			view := args[0]
			name := args[1]

			// Get the HTTP client
			client := NewApiClient(cmd)
			notifier, err := client.Notifiers().Get(view, name)

			if err != nil {
				return fmt.Errorf("Error fetching notifier: %s", err)
			}

			var output []string
			output = append(output, "Name | EntityType")
			output = append(output, fmt.Sprintf("%v | %v", notifier.Name, notifier.Entity))

			printTable(cmd, output)

			return nil
		},
	}

	return &cmd
}
