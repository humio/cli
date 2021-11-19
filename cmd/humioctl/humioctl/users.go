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
	"github.com/humio/cli/api"
	format2 "github.com/humio/cli/cmd/humioctl/internal/format"
	"github.com/spf13/cobra"
)

func newUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users [Root Only]",
	}

	cmd.AddCommand(newUsersAddCmd())
	cmd.AddCommand(newUsersRemoveCmd())
	cmd.AddCommand(newUsersUpdateCmd())
	cmd.AddCommand(newUsersListCmd())
	cmd.AddCommand(newUsersShowCmd())
	cmd.AddCommand(newUsersRotateApiTokenCmd())

	return cmd
}

func printUserDetailsTable(cmd *cobra.Command, user api.User) {
	details := [][]format2.Value{
		{format2.String("Username"), format2.String(user.Username)},
		{format2.String("Name"), format2.String(user.FullName)},
		{format2.String("Is Root"), format2.YesNo(user.IsRoot)},
		{format2.String("Email"), format2.String(user.Email)},
		{format2.String("Created At"), format2.String(user.CreatedAt)},
		{format2.String("Country Code"), format2.String(user.CountryCode)},
		{format2.String("Company"), format2.String(user.Company)},
		{format2.String("ID"), format2.String(user.ID)},
	}

	format2.PrintDetailsTable(cmd, details)
}
