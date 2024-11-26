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
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
)

func installPackageCmd() *cobra.Command {
	var (
		ownership string
	)

	cmd := cobra.Command{
		Use:   "install <repo-or-view> <path-to-package-dir>",
		Short: "Installs a package.",
		Long: `
Packages can be installed from a directory, Github Repository URL, Zip File, or
Zip File URL.

  $ humioctl packages install myrepo /path/to/package/dir/ --ownership organization 
  $ humioctl packages install myrepo /path/to/pazkage.zip --ownership organization
  $ humioctl packages install myrepo https://github.com/org/mypackage-name -ownership organization
  $ humioctl packages install myrepo https://content.example.com/mypackage-name.zip -ownership organization
`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			path := args[1]
			client := NewApiClient(cmd)

			if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
				downloadedFile, err := getURLPackage(path)
				exitOnError(cmd, err, fmt.Sprintf("Failed to download file to path: %s", path))

				path = downloadedFile.Name()
			}

			isDir, err := isDirectory(path)
			exitOnError(cmd, err, "Errors installing archive")

			if ownership == "" {
				ownership = "user"
			}

			if ownership != "organization" && ownership != "user" {
				cmd.PrintErrln("query ownership must be set to either `organization` or `user`")
				os.Exit(1)
			}

			var validationResult *api.ValidationResponse
			if isDir {
				validationResult, err = client.Packages().InstallFromDirectory(path, repoOrViewName, ownership)
			} else {
				validationResult, err = client.Packages().InstallArchive(repoOrViewName, path, ownership)
			}
			exitOnError(cmd, err, "Errors installing archive")

			if !validationResult.IsValid() {
				printValidation(cmd, validationResult)
				os.Exit(1)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Successfully installed package in: %s\n", path)
		},
	}

	cmd.Flags().StringVarP(&ownership, "ownership", "o", "", "The query ownership of installed queries e.g. in triggers. Possible are `organization` and `user`. Defaults to `user`")

	return &cmd
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func getURLPackage(url string) (*os.File, error) {
	// Github URLs should have the .git suffix removed
	if strings.HasPrefix(url, "https://github.com") || strings.HasPrefix(url, "http://github.com") {
		url = strings.TrimSuffix(url, ".git")
	}

	zipBallURL := url + "/zipball/master/"
	// #nosec G107
	response, err := http.Get(zipBallURL)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, fmt.Errorf("failed to get response")
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("error downloading file %s: %s", zipBallURL, response.Status)
	}

	defer func() {
		_ = response.Body.Close()
	}()
	zipContent, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	tempDir := os.TempDir()
	zipFile, err := os.CreateTemp(tempDir, "humio-package.*.zip")
	if err != nil {
		return nil, err
	}

	_, err = zipFile.Write(zipContent)
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}
