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
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

const viewTypeName = "View"

func newViewsListCmd() *cobra.Command {
	var viewOnly bool

	cmd := cobra.Command{
		Use:   "list",
		Short: "Lists all views you have access to",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			views, err := client.Views().List()
			exitOnError(cmd, err, "Error while fetching view list")

			rows := make([][]format.Value, len(views))
			for i, view := range views {
				if viewOnly {
					if view.Typename == viewTypeName {
						rows[i] = []format.Value{format.String(view.Name)}
					}
				} else {
					rows[i] = []format.Value{format.String(view.Name)}
				}
			}

			printOverviewTable(cmd, []string{"Name"}, rows)
		},
	}

	cmd.Flags().BoolVar(&viewOnly, "only-views", false, "Display only Views (i.e. do not include Repositories).")

	return &cmd
}
