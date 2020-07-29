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
	"github.com/spf13/cobra"
)

func newUsersListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all users. [Root Only]",
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			users, err := client.Users().List()
			exitOnError(cmd, err, "error fetching user list")

			rows := make([]string, len(users))
			for i, user := range users {
				rows[i] = formatSimpleAccount(user)
			}

			printTable(cmd, append([]string{
				"Username | Name | Root | Created"},
				rows...,
			))
		},
	}
}
