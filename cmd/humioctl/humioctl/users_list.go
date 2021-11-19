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
	"github.com/spf13/cobra"
)

func newUsersListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all users. [Root Only]",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			users, err := client.Users().List()
			helpers.ExitOnError(cmd, err, "Error fetching user list")

			rows := make([][]format2.Value, len(users))
			for i, user := range users {
				rows[i] = []format2.Value{
					format2.String(user.Username),
					format2.String(user.FullName),
					format2.YesNo(user.IsRoot),
					format2.String(user.CreatedAt),
					format2.String(user.ID),
				}
			}

			format2.PrintOverviewTable(cmd, []string{"Username", "Name", "Root", "Created", "ID"}, rows)
		},
	}
}
