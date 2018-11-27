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
	"strings"

	"github.com/humio/cli/api"
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

	return cmd
}

func formatSimpleAccount(account api.User) string {
	columns := []string{account.Username, account.FullName, yesNo(account.IsRoot), account.CreatedAt}
	return strings.Join(columns, " | ")
}

func printUserTable(user api.User) {
	userData := []string{user.Username, user.FullName, user.CreatedAt, yesNo(user.IsRoot)}

	printTable([]string{
		"Username | Name | Created At | Is Root",
		strings.Join(userData, "|"),
	})
}
