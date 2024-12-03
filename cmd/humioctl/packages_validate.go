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

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
)

func validatePackageCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "validate <repo-or-view> <package-dir-or-zip>",
		Short: "Validate a package's content.",
		Long: `
Packages can be validated from a directory or Zip File. You must specify the
repository or view to validate the package against.

  $ humioctl packages validate myrepo /path/to/package/dir/
  $ humioctl packages validate myrepo /path/to/pazkage.zip
`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			dirPath := args[1]
			client := NewApiClient(cmd)

			if !filepath.IsAbs(dirPath) {
				var err error
				dirPath, err = filepath.Abs(dirPath)
				exitOnError(cmd, err, "Invalid path")
				dirPath += "/"
			}

			validationResult, err := client.Packages().Validate(repoOrViewName, dirPath)
			exitOnError(cmd, err, "Errors validating package")

			if !validationResult.IsValid() {
				printValidation(cmd, validationResult)
				os.Exit(1)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Successfully validated package in: %s\n", dirPath)
		},
	}

	return &cmd
}

func printValidation(cmd *cobra.Command, validationResult *api.ValidationResponse) {
	cmd.PrintErrln("Package is not valid")
	cmd.PrintErrln(validationResult.InstallationErrors)
	cmd.PrintErrln(validationResult.ParseErrors)
}
