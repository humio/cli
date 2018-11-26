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

	"github.com/spf13/cobra"
)

// listCmd represents the list command
func newUsersListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewApiClient(cmd)

			users, err := client.Users().List()

			if err != nil {
				return fmt.Errorf("error fetching user list: %s", err)
			}

			rows := make([]string, len(users))
			for i, user := range users {
				rows[i] = formatSimpleAccount(user)
			}

			printTable(append([]string{
				"Username | Name | Root | Created"},
				rows...,
			))

			return nil
		},
	}
}
