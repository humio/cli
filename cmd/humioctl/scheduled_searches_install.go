// Copyright © 2024 CrowdStrike
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

	"github.com/humio/cli/internal/api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newScheduledSearchesInstallCmd() *cobra.Command {
	var filePath, url, name string

	cmd := cobra.Command{
		Use:   "install [flags] <view>",
		Short: "Installs a scheduled search in a view",
		Long: `Install a scheduled search from a URL or from a local file.

The install command allows you to install scheduled searches from a URL or from a local file, e.g.

  $ humioctl scheduled-searches install viewName --url=https://example.com/acme/scheduled-search.yaml

  $ humioctl scheduled-searches install viewName --file=./scheduled-searches.yaml
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var content []byte
			var err error

			// Check that we got the right number of argument
			// if we only got <view> you must supply --file or --url.
			if l := len(args); l == 1 {
				if filePath != "" {
					content, err = getBytesFromFile(filePath)
				} else if url != "" {
					content, err = getBytesFromURL(url)
				} else {
					cmd.Printf("You must specify a path using --file or --url\n")
					os.Exit(1)
				}
			}
			exitOnError(cmd, err, "Could to load the scheduled search")

			client := NewApiClient(cmd)
			viewName := args[0]

			var scheduledSearch api.ScheduledSearch
			err = yaml.Unmarshal(content, &scheduledSearch)
			exitOnError(cmd, err, "Could not unmarshal the scheduled search")

			if name != "" {
				scheduledSearch.Name = name
			}

			_, err = client.ScheduledSearches().Create(viewName, &scheduledSearch)
			exitOnError(cmd, err, "Could not create the scheduled search")

			fmt.Fprintln(cmd.OutOrStdout(), "Scheduled search created")
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "The local file path to the scheduled search to install.")
	cmd.Flags().StringVar(&url, "url", "", "A URL to fetch the scheduled search file from.")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Install the scheduled search under a specific name, ignoring the `name` attribute in the scheduled search file.")
	cmd.MarkFlagsMutuallyExclusive("file", "url")
	return &cmd
}
