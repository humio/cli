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
	"github.com/spf13/cobra"
)

func newViewsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <view> \"descriptive reason for why it is being deleted\"",
		Short: "Delete a view.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			view := args[0]
			reason := args[1]

			fmt.Printf("Deleting view %s with reason %q\n", view, reason)

			client := NewApiClient(cmd)

			apiError := client.Views().Delete(view, reason)
			exitOnError(cmd, apiError, "error removing view")
		},
	}
}
