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
	"os"

	"github.com/humio/cli/prompt"

	"github.com/spf13/cobra"
)

func installPackageCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "install [flags] <repo-or-view-name> <path-to-package-dir>",
		Short: "Installs a package.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			out := prompt.NewPrompt(cmd.OutOrStdout())
			repoOrView := args[0]
			path := args[1]

			out.Info(fmt.Sprintf("Installing Package from: %s", path))

			// Get the HTTP client
			client := NewApiClient(cmd)

			var createErr error
			isDir, err := isDirectory(path)
			println(path)
			if err != nil {
				out.Error(fmt.Sprintf("Errors installing archive: %s", err))
				os.Exit(1)
			}

			if isDir {
				createErr = client.Packages().InstallFromDirectory(path, repoOrView)
			} else {
				createErr = client.Packages().InstallArchive(repoOrView, path)
			}
			if createErr != nil {
				out.Error(fmt.Sprintf("Errors installing archive: %s", createErr))
				os.Exit(1)
			}
		},
	}

	return &cmd
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
