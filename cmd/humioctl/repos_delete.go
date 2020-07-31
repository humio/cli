// Copyright © 2020 Humio Ltd.
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

func newReposDeleteCmd() *cobra.Command {
	var allowDataDeletionFlag bool

	cmd := cobra.Command{
		Use:   "delete [flags] <repo> \"descriptive reason for why it is being deleted\"",
		Short: "Delete a repository.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			reason := args[1]

			client := NewApiClient(cmd)

			apiError := client.Repositories().Delete(repo, reason, allowDataDeletionFlag)
			exitOnError(cmd, apiError, "error removing repository")
		},
	}

	cmd.Flags().BoolVar(&allowDataDeletionFlag, "allow-data-deletion", false, "Allow changing retention settings for a non-empty repository")

	return &cmd
}
