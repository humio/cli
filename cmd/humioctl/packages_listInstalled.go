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

func listInstalledPackagesCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list-installed [flags] <view-or-repo-name>",
		Short: "List all installed packages in a repository.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			viewName := args[0]

			// Get the HTTP client
			client := NewApiClient(cmd)
			installedPackages, err := client.Packages().ListInstalled(viewName)

			if err != nil {
				return fmt.Errorf("Error fetching packages: %s", err)
			}

			var output []string
			output = append(output, "ID | InstalledBy | UpdatedBy | Source | AvailableUpdate")
			for i := 0; i < len(installedPackages); i++ {
				installedPackage := installedPackages[i]
				output = append(output, fmt.Sprintf("%v | %v | %v | %v | %v", installedPackage.ID, installedPackage.InstalledBy.Username, valueOrEmpty(installedPackage.UpdatedBy.Username), installedPackage.Source, valueOrEmpty(installedPackage.AvailableUpdate)))
			}

			printTable(cmd, output)

			return nil
		},
	}

	return &cmd
}
