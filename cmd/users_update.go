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

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
)

func newUsersUpdateCmd() *cobra.Command {
	var rootFlag boolPtrFlag
	var nameFlag stringPtrFlag

	cmd := cobra.Command{
		Use:   "update",
		Short: "Updates a user's settings and global permissions.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			username := args[0]

			client := NewApiClient(cmd)

			user, err := client.Users().Update(username, api.UserChangeSet{
				IsRoot:   rootFlag.value,
				FullName: nameFlag.value,
			})

			if err != nil {
				return fmt.Errorf("Error updating user: %s", err)
			}

			printUserTable(user)

			return nil
		},
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")
	cmd.Flags().Var(&rootFlag, "root", "If true grants root access to the user.")
	cmd.Flags().Var(&nameFlag, "name", "The full name of the user.")
	// updateCmd.Flags().VarP("name", false, "Sets the full name of the user.")
	// updateCmd.Flags().VarP("email", false, "Sets the email of the user (this will not change the username if email is used).")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return &cmd
}
