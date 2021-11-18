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
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newUsersRemoveCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "remove",
		Short: "Remove a user [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			client := NewApiClient(cmd)

			removedUser, err := client.Users().Remove(username)
			helpers.ExitOnError(cmd, err, "Error removing the user")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed user with username %q with ID %q\n", removedUser.Username, removedUser.ID)
		},
	}

	return &cmd
}
