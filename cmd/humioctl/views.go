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
	"os"

	"github.com/humio/cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func newViewsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "views",
		Short: "Manage views",
	}

	cmd.AddCommand(newViewsShowCmd())
	cmd.AddCommand(newViewsListCmd())

	return cmd
}

func printViewTable(view *api.View) {

	data := [][]string{
		{"Name", view.Name},
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.AppendBulk(data)
	w.SetBorder(false)
	w.SetColumnSeparator(":")
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	fmt.Println()
	w.Render()
	fmt.Println()
}

func printViewRoleTable(view *api.View) {

	data := [][]string{}

	for _, role := range view.Roles {
		data = append(data, []string{role.Role.Name, role.QueryPrefix})
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.AppendBulk(data)
	w.SetBorder(true)
	w.SetHeader([]string{"Role", "Query Prefix"})
	w.SetColumnSeparator(":")
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	fmt.Println()
	w.Render()
	fmt.Println()
}

func printViewConnectionsTable(view *api.View) {
	if len(view.Connections) == 0 {
		return
	}

	data := [][]string{}

	for _, conn := range view.Connections {
		data = append(data, []string{conn.RepoName, conn.Filter})
	}

	w := tablewriter.NewWriter(os.Stdout)
	w.AppendBulk(data)
	w.SetBorder(true)
	w.SetHeader([]string{"Repository", "Query Prefix"})
	w.SetColumnSeparator(":")
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	fmt.Println()
	w.Render()
	fmt.Println()
}
