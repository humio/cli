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
	"github.com/humio/cli/internal/api"
	"github.com/humio/cli/internal/format"
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

func printUserDetailsTable(cmd *cobra.Command, user api.User) {
	details := [][]format.Value{
		{format.String("Username"), format.String(user.Username)},
		{format.String("Name"), format.StringPtr(user.FullName)},
		{format.String("Is Root"), format.Bool(user.IsRoot)},
		{format.String("Email"), format.StringPtr(user.Email)},
		{format.String("Created At"), format.String(user.CreatedAt)},
		{format.String("Country Code"), format.StringPtr(user.CountryCode)},
		{format.String("Company"), format.StringPtr(user.Company)},
		{format.String("ID"), format.String(user.ID)},
	}

	printDetailsTable(cmd, details)
}
