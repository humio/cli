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
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
)

func newViewsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "views",
		Short: "Manage views",
	}

	cmd.AddCommand(newViewsShowCmd())
	cmd.AddCommand(newViewsListCmd())
	cmd.AddCommand(newViewsCreateCmd())
	cmd.AddCommand(newViewsUpdateCmd())
	cmd.AddCommand(newViewsDeleteCmd())

	return cmd
}

func printViewConnectionsTable(cmd *cobra.Command, view *api.View) {
	if len(view.Connections) == 0 {
		return
	}

	var rows [][]string
	for _, conn := range view.Connections {
		rows = append(rows, []string{view.Name, conn.RepoName, conn.Filter})
	}

	printOverviewTable(cmd, []string{"View", "Repository", "Query Prefix"}, rows)
}
