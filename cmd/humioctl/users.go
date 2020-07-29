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
	"strings"

	"github.com/humio/cli/api"
	"github.com/olekukonko/tablewriter"
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

func printUserTable(cmd *cobra.Command, user api.User) {

	data := [][]string{
		[]string{"Username", user.Username},
		[]string{"Name", user.FullName},
		[]string{"Is Root", yesNo(user.IsRoot)},
		[]string{"Email", user.Email},
		[]string{"Created At", user.CreatedAt},
		[]string{"Country Code", user.CountryCode},
		[]string{"Company", user.Company},
	}

	w := tablewriter.NewWriter(cmd.OutOrStdout())
	w.AppendBulk(data)
	w.SetBorder(false)
	w.SetColumnSeparator(":")
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	fmt.Println()
	w.Render()
	fmt.Println()
}
