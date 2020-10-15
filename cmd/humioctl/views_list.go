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
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func newViewsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all views you have access to",
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)

			views, apiErr := client.Views().List()
			if apiErr != nil {
				return nil, fmt.Errorf("error while fetching view list: %w", apiErr)
			}

			rows := make([][]string, len(views))
			for i, view := range views {
				rows[i] = []string{view.Name}
			}

			w := tablewriter.NewWriter(cmd.OutOrStdout())
			w.AppendBulk(rows)
			w.SetBorder(false)

			cmd.Println()
			w.Render()
			cmd.Println()

			return nil, nil
		}),
	}
}
