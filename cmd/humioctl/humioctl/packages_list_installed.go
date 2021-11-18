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

package humioctl

import (
	format2 "github.com/humio/cli/cmd/humioctl/internal/format"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
)

func listInstalledPackagesCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list-installed <repo-or-view>",
		Short: "List all installed packages in a repository.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			client := NewApiClient(cmd)

			installedPackages, err := client.Packages().ListInstalled(repoOrViewName)
			helpers.ExitOnError(cmd, err, "Error fetching packages")

			var rows [][]format2.Value
			for _, installedPackage := range installedPackages {
				rows = append(rows, []format2.Value{
					format2.String(installedPackage.ID),
					format2.String(installedPackage.InstalledBy.Username),
					format2.ValueOrEmpty(installedPackage.UpdatedBy.Username),
					format2.String(installedPackage.Source),
					format2.ValueOrEmpty(installedPackage.AvailableUpdate),
				})
			}

			format2.PrintOverviewTable(cmd, []string{"ID", "Installed By", "Updated By", "Source", "Available Update"}, rows)
		},
	}

	return &cmd
}
