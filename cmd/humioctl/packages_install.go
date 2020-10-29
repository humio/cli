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
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/humio/cli/api"
	"github.com/humio/cli/prompt"

	"github.com/spf13/cobra"
)

func installPackageCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "install [flags] <repo-or-view-name> <path-to-package-dir>",
		Short: "Installs a package.",
		Long: `
Packages can be installed from a directory, Github Repository URL, Zip File, or
Zip File URL.

  $ humioctl packages install myrepo /path/to/package/dir/
  $ humioctl packages install myrepo /path/to/pazkage.zip
  $ humioctl packages install myrepo https://github.com/org/mypackage-name
  $ humioctl packages install myrepo https://content.example.com/mypackage-name.zip
`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			out := prompt.NewPrompt(cmd.OutOrStdout())
			repoOrView := args[0]
			path := args[1]

			out.Info(fmt.Sprintf("Installing Package from: %s", path))

			if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
				downloadedFile, err := getURLPackage(path)

				if err != nil {
					out.Error(fmt.Sprintf("Failed to download: %s %s", path, err))
					os.Exit(1)
				}

				// defer os.Remove(downloadedFile.Name())

				path = downloadedFile.Name()
			}

			isDir, err := isDirectory(path)

			if err != nil {
				out.Error(fmt.Sprintf("Errors installing archive: %s", err))
				os.Exit(1)
			}

			// Get the HTTP client
			client := NewApiClient(cmd)

			var validationResult *api.ValidationResponse
			var createErr error

			if isDir {
				validationResult, createErr = client.Packages().InstallFromDirectory(path, repoOrView)
			} else {
				validationResult, createErr = client.Packages().InstallArchive(repoOrView, path)
			}

			if createErr != nil {
				out.Error(fmt.Sprintf("Errors installing archive: %s", createErr))
				os.Exit(1)
			} else if !validationResult.IsValid() {
				printValidation(out, validationResult)
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

func getURLPackage(url string) (*os.File, error) {
	// Github URLs should have the .got suffix removed
	if strings.HasPrefix(url, "https://github.com") || strings.HasPrefix(url, "http://github.com") {
		url = strings.TrimSuffix(url, ".git")
	}

	zipBallURL := url + "/zipball/master/"
	response, err := http.Get(zipBallURL)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("error downloading file %s: %s", zipBallURL, response.Status)
	}

	defer response.Body.Close()
	zipContent, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	tempDir := os.TempDir()
	zipFile, err := ioutil.TempFile(tempDir, "humio-package.*.zip")

	if err != nil {
		return nil, err
	}

	_, err = zipFile.Write(zipContent)

	if err != nil {
		return nil, err
	}

	return zipFile, nil
}
