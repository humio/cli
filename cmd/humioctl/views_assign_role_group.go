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

func newViewsAssignRoleGroupCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "assign <role> <group> <view>",
		Short: "Assign Role to a Group for a View",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			roleName := args[0]
			groupName := args[1]
			viewName := args[2]
			client := NewApiClient(cmd)

			role, err := client.Roles().Get(roleName)
			exitOnError(cmd, err, "Error fetching role")
			
			group, err := client.Groups().Get(groupName)
			exitOnError(cmd, err, "Error fetching group")

			err = client.Views().AssignRoleToGroup(viewName, group.ID, role.ID)
			exitOnError(cmd, err, "Error assigning permission")
		},
	}

	return &cmd
}
