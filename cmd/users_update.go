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
	var root boolPtrFlag

	cmd := cobra.Command{
		Use:   "update",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		ValidArgs: []string{"username"},
		RunE: func(cmd *cobra.Command, args []string) error {

			username := args[0]

			client := NewApiClient(cmd)

			user, err := client.Users().Update(username, api.UserChangeSet{IsRoot: root.value})

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
	cmd.Flags().Var(&root, "root", "If true grants root access to the user.")
	// updateCmd.Flags().VarP("name", false, "Sets the full name of the user.")
	// updateCmd.Flags().VarP("email", false, "Sets the email of the user (this will not change the username if email is used).")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return &cmd
}
