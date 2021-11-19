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

package humioctl

import (
	format2 "github.com/humio/cli/cmd/humioctl/internal/format"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"sort"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
)

func newReposListCmd() *cobra.Command {
	var orderBySize, reverse bool

	cmd := cobra.Command{
		Use:   "list [flags]",
		Short: "List repositories.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			repos, err := client.Repositories().List()
			helpers.ExitOnError(cmd, err, "Error fetching repository")

			sort.Slice(repos, func(i, j int) bool {
				var a, b api.RepoListItem
				if reverse {
					a = repos[i]
					b = repos[j]
				} else {
					a = repos[j]
					b = repos[i]
				}

				if orderBySize {
					return a.SpaceUsed > b.SpaceUsed
				}
				return a.Name < b.Name
			})

			rows := make([][]format2.Value, len(repos))
			for i, view := range repos {
				rows[i] = []format2.Value{
					format2.String(view.Name),
					format2.ByteCountDecimal(view.SpaceUsed),
					format2.String(view.ID),
				}
			}

			format2.PrintOverviewTable(cmd, []string{"Name", "Space Used", "ID"}, rows)
		},
	}

	cmd.Flags().BoolVarP(&orderBySize, "size", "s", false, "Order by size instead of name")
	cmd.Flags().BoolVarP(&reverse, "reverse", "r", true, "Reverse sorting order")

	return &cmd
}
