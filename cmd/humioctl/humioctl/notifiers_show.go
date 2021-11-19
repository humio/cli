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
	format2 "github.com/humio/cli/cmd/humioctl/internal/format"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func newNotifiersShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show <repo-or-view> <name>",
		Short: "Show details about a notifier in a repository or view.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			name := args[1]
			client := NewApiClient(cmd)

			notifier, err := client.Notifiers().Get(repoOrViewName, name)
			helpers.ExitOnError(cmd, err, "Error fetching notifier")

			details := [][]format2.Value{
				{format2.String("Name"), format2.String(notifier.Name)},
				{format2.String("EntityType"), format2.String(notifier.Entity)},
			}

			format2.PrintDetailsTable(cmd, details)
		},
	}

	return &cmd
}
