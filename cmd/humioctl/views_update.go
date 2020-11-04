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

func newViewsUpdateCmd() *cobra.Command {
	connections := make(map[string] string)

	cmd := cobra.Command{
		Use:   "update",
		Short: "Updates the settings of a view",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]

			if len(connections) == 0 {
				exitOnError(cmd, fmt.Errorf("you must specify at least one connection flag"), "nothing specified to update")
			}

			client := NewApiClient(cmd)

			err := client.Views().UpdateConnections(viewName, connections)
			exitOnError(cmd, err, "error updating view connections")

			view, apiErr := client.Views().Get(viewName)
			exitOnError(cmd, apiErr, "error fetching view")
			printViewTable(view)
			fmt.Println()
		},
	}

	cmd.Flags().StringToStringVar(&connections, "connection", connections, "Sets a repository connection with the chosen filter.")

	return &cmd
}
