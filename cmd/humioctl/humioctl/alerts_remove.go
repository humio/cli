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

package humioctl

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newAlertsRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [flags] <view> <name>",
		Short: "Removes an alert.",
		Long:  `Removes the alert with name '<name>' in the view with name '<view>'.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			viewName := args[0]
			alertName := args[1]
			client := NewApiClient(cmd)

			err := client.Alerts().Delete(viewName, alertName)
			helpers.ExitOnError(cmd, err, "Error removing alert")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed alert %q from view %q\n", alertName, viewName)
		},
	}

	return cmd
}
