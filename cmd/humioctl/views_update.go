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

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
)

func newViewsUpdateCmd() *cobra.Command {
	var enableAutomaticSearchFlag, disableAutomaticSearchFlag bool

	connsFlag := []string{}
	connections := []api.ViewConnectionInput{}
	description := ""

	cmd := cobra.Command{
		Use:   "update [flags] <view>",
		Short: "Updates the settings of a view",
		Long: `Updates the settings of a view with the provided arguments.

The "description" flag can be specified multiple times for adding multiple connections to the view.
If you want to query all events you can specify a wildcard as the filter.

Here's an example that updates a view named "important-view" to search all data in the two repositories,
namely "repo1" and "repo2":

  $ humioctl views update important-view --connection "repo1=*" --connection="repo2=*" --description "very important view"
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]
			client := NewApiClient(cmd)

			if len(connsFlag) == 0 && description == "" && !enableAutomaticSearchFlag && !disableAutomaticSearchFlag {
				exitOnError(cmd, fmt.Errorf("you must specify at least one flag"), "Nothing specified to update")
			}

			if len(connsFlag) > 0 {
				for _, v := range connsFlag {
					parts := strings.SplitN(v, "=", 2)
					if len(parts) != 2 {
						exitOnError(cmd, fmt.Errorf("all connections must follow the format: <repoName>=<filterString>"), "Error updating view connections")
					}

					repo := parts[0]
					filter := parts[1]

					connections = append(
						connections,
						api.ViewConnectionInput{
							RepositoryName: repo,
							Filter:         filter,
						})
				}
				err := client.Views().UpdateConnections(viewName, connections)
				exitOnError(cmd, err, "Error updating view connections")
			}

			if description != "" {
				err := client.Views().UpdateDescription(viewName, description)
				exitOnError(cmd, err, "Error updating view description")
			}

			if disableAutomaticSearchFlag {
				err := client.Views().UpdateAutomaticSearch(viewName, false)
				exitOnError(cmd, err, "Error disabling automatic search")
			}
			if enableAutomaticSearchFlag {
				err := client.Views().UpdateAutomaticSearch(viewName, true)
				exitOnError(cmd, err, "Error enabling automatic search")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully updated view %q\n", viewName)
		},
	}

	cmd.Flags().StringArrayVar(&connsFlag, "connection", connsFlag, "Sets a repository connection with the chosen filter in format: <repoName>=<filterString>")
	cmd.Flags().StringVar(&description, "description", description, "Sets the view description.")
	cmd.Flags().BoolVar(&enableAutomaticSearchFlag, "enable-automatic-search", false, "Enable automatic search for the view.")
	cmd.Flags().BoolVar(&disableAutomaticSearchFlag, "disable-automatic-search", false, "Disable automatic search for the view.")
	cmd.MarkFlagsMutuallyExclusive("enable-automatic-search", "disable-automatic-search")

	return &cmd
}
