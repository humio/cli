// Copyright © 2018 Humio Ltd.
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

func newViewsCreateCmd() *cobra.Command {
	connsFlag := make(map[string]string)
	connections := make(map[string][]string)
	description := ""

	cmd := &cobra.Command{
		Use:   "create [flags] <view-name>",
		Short: "Create a view.",
		Long: `Creates a view with the provided arguments.

The "description" flag is a string, and the "connections" flag is a comma-separated list of key-value pairs
where the key is the repository name and the value being the filter applied to the queries in that repository.
If you want to query all events you can specify a wildcard as the filter.

Here's an example that creates a view named "important-view" to search all data in the two repositories,
namely "repo1" and "repo2":

  $ humioctl views create important-view --connection "repo1=*,repo2=*" --description "very important view"
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]
			client := NewApiClient(cmd)

			if len(connsFlag) == 0 {
				exitOnError(cmd, fmt.Errorf("you must specify at least view connection"), "Error creating view")
			}

			for k, v := range connsFlag {
				connections[k] = append(connections[k], v)
			}

			err := client.Views().Create(viewName, description, connections)
			exitOnError(cmd, err, "Error creating view")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created view: %q\n", viewName)
		},
	}

	cmd.Flags().StringToStringVar(&connsFlag, "connection", connsFlag, "Sets a repository connection with the chosen filter.")
	cmd.Flags().StringVar(&description, "description", description, "Sets an optional description")

	return cmd
}
