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

	"github.com/humio/cli/api"

	"github.com/humio/cli/prompt"

	"github.com/spf13/cobra"
)

func validatePackageCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "validate [flags] <repo-or-view-name> <package-dir>",
		Short: "Validate a package's content.",
		Long: `
Packages can be validated from a directory or Zip File. You must specify the
repository or view to validate the package against.

  $ humioctl packages validate myrepo /path/to/package/dir/
  $ humioctl packages validate myrepo /path/to/pazkage.zip
`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			out := prompt.NewPrompt(cmd.OutOrStdout())

			viewName := args[0]
			dirPath := args[1]

			if !filepath.IsAbs(dirPath) {
				var err error
				dirPath, err = filepath.Abs(dirPath)
				if err != nil {
					out.Error(fmt.Sprintf("Invalid path: %s", err))
					os.Exit(1)
				}
				dirPath += "/"
			}

			out.Info(fmt.Sprintf("Validating Package in: %s", dirPath))

			// Get the HTTP client
			client := NewApiClient(cmd)

			validationResult, apiErr := client.Packages().Validate(viewName, dirPath)
			if apiErr != nil {
				out.Error(fmt.Sprintf("Errors validating package: %s", apiErr))
				os.Exit(1)
			}

			if validationResult.IsValid() {
				out.Info("Package is valid")
			} else {
				printValidation(out, validationResult)
				os.Exit(1)
			}
		},
	}

	return &cmd
}

func printValidation(out *prompt.Prompt, validationResult *api.ValidationResponse) {
	out.Error("Package is not valid")
	out.Error(out.List(validationResult.InstallationErrors))
	out.Error(out.List(validationResult.ParseErrors))
}
