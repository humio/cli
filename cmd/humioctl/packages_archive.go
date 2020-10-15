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
	"path/filepath"

	"github.com/humio/cli/prompt"

	"github.com/spf13/cobra"
)

func archivePackageCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "archive [flags] <package-dir> <output-file>",
		Short: " Create a zip containing the content of a package directory.",
		Args:  cobra.ExactArgs(2),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			out := prompt.NewPrompt(cmd.OutOrStdout())
			dirPath := args[0]
			outPath := args[1]

			if !filepath.IsAbs(dirPath) {
				var err error
				dirPath, err = filepath.Abs(dirPath)
				if err != nil {
					return nil, fmt.Errorf("invalid path: %w", err)
					os.Exit(1)
				}
				dirPath += "/"
			}

			out.Info(fmt.Sprintf("Archiving Package in: %s", dirPath))

			// Get the HTTP client
			client := NewApiClient(cmd)

			createErr := client.Packages().CreateArchive(dirPath, outPath)
			if createErr != nil {
				return nil, fmt.Errorf("errors creating archive: %w", createErr)
			}

			return fmt.Sprintf("Created %s with package contents from %s.", outPath, dirPath), nil
		}),
	}

	return &cmd
}
