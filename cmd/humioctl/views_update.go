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
	connections := make(map[string]string)
	description := ""

	cmd := cobra.Command{
		Use:   "update [flags] <view>",
		Short: "Updates the settings of a view",
		Long: `Updates the settings of a view with the provided arguments.

The "description" flag is a string, and the "connections" flag is a comma-separated list of key-value pairs
where the key is the repository name and the value being the filter applied to the queries in that repository.
If you want to query all events you can specify a wildcard as the filter.

Here's an example that updates a view named "important-view" to search all data in the two repositories,
namely "repo1" and "repo2":

  $ humioctl views update important-view --connection "repo1=*,repo2=*" --description "very important view"
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]
			client := NewApiClient(cmd)

			if len(connections) == 0 && description == "" {
				exitOnError(cmd, fmt.Errorf("you must specify at least one flag"), "Nothing specified to update")
			}

			if len(connections) > 0 {
				err := client.Views().UpdateConnections(viewName, connections)
				exitOnError(cmd, err, "Error updating view connections")
			}

			if description != "" {
				err := client.Views().UpdateDescription(viewName, description)
				exitOnError(cmd, err, "Error updating view description")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully updated view %q\n", viewName)
		},
	}

	cmd.Flags().StringToStringVar(&connections, "connection", connections, "Sets a repository connection with the chosen filter.")
	cmd.Flags().StringVar(&description, "description", description, "Sets the view description.")

	return &cmd
}
