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

package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
func newViewsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all views you have access to",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			views, err := client.Views().List()

			if err != nil {
				cmd.Println(fmt.Errorf("error fetching view list: %s", err))
				os.Exit(1)
			}

			rows := make([][]string, len(views))
			for i, view := range views {
				rows[i] = []string{view.Name}
			}

			w := tablewriter.NewWriter(cmd.OutOrStderr())
			w.AppendBulk(rows)
			w.SetBorder(false)

			fmt.Println()
			w.Render()
			fmt.Println()
		},
	}
}
