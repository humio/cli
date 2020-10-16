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

	"github.com/humio/cli/prompt"

	"github.com/spf13/cobra"
)

func listInstalledPackagesCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "listInstalled [flags] <view-or-repo-name>",
		Short: "List all installed packages in a repository.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := prompt.NewPrompt(cmd.OutOrStdout())
			viewName := args[0]

			out.Info(fmt.Sprintf("Listing installed packages in view %s", viewName))

			// Get the HTTP client
			client := NewApiClient(cmd)

			installedPackages, err := client.Packages().ListInstalled(viewName)

			if err != nil {
				return fmt.Errorf("Error fetching parsers: %s", err)
			}

			for i := 0; i < len(installedPackages); i++ {
				installedPackage := installedPackages[i]
				out.Info(fmt.Sprintf(installedPackage.ID))
			}

			return nil
		},
	}

	return &cmd
}
