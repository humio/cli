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
	"os"

	"github.com/spf13/cobra"
)

func newNotifiersRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <repo-or-view> <name>",
		Short: "Removes an alert notifier.",
		Long:  `Removes the alert notifier with name '<name>' in the view with name '<repo-or-view>'.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			name := args[1]
			client := NewApiClient(cmd)

			err := client.Notifiers().Delete(repoOrViewName, name)
			if err != nil {
				cmd.Printf("Error removing notifier: %s\n", err)
				os.Exit(1)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Notifier removed")
		},
	}

	return cmd
}
