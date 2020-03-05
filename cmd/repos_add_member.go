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
	"github.com/spf13/cobra"
)

func newReposAddMemberCmd() *cobra.Command {
	var adminRights, deleteRights bool

	cmd := cobra.Command{
		Use:   "add-member [flags] <repo> <username>",
		Short: "Adds a user to a repository.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			userName := args[1]

			client := NewApiClient(cmd)

			_, apiErr := client.Repositories().AddMember(repoName, userName, adminRights, deleteRights)
			exitOnError(cmd, apiErr, "error creating repository")
		},
	}

	cmd.Flags().BoolVar(&adminRights, "admin-rights", false, "Grant the user admin rights")
	cmd.Flags().BoolVar(&deleteRights, "delete-rights", false, "Grant the user delete rights")

	return &cmd
}
