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

func uninstallPackageCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "uninstall <repo-or-view> <package>",
		Short: "Uninstall a package.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)
			repoOrViewName := args[0]
			packageName := args[1]

			err := client.Packages().UninstallPackage(repoOrViewName, packageName)
			helpers.ExitOnError(cmd, err, "Errors uninstalling package")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully uninstalled package %s from view/repo %s\n", packageName, repoOrViewName)
		},
	}

	return &cmd
}
