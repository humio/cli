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

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newAggregateAlertsInstallCmd() *cobra.Command {
	var (
		filePath, url string
	)

	cmd := cobra.Command{
		Use:   "install [flags] <view>",
		Short: "Installs an aggregate alert in a view",
		Long: `Install an aggregate alert from a URL or from a local file.

The install command allows you to install aggregate alerts from a URL or from a local file, e.g.

  $ humioctl aggregate-alerts install viewName --url=https://example.com/acme/aggregate-alert.yaml

  $ humioctl aggregate-alerts install viewName --file=./aggregate-alert.yaml
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
			exitOnError(cmd, err, "Could to load the aggregate alert")

			client := NewApiClient(cmd)
			viewName := args[0]

			var aggregateAlert api.AggregateAlert
			err = yaml.Unmarshal(content, &aggregateAlert)
			exitOnError(cmd, err, "Could not unmarshal the aggregate alert")

			_, err = client.AggregateAlerts().Create(viewName, &aggregateAlert)
			exitOnError(cmd, err, "Could not create the aggregate alert")

			fmt.Fprintln(cmd.OutOrStdout(), "Aggregate alert created")
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "The local file path to the aggregate alert to install.")
	cmd.Flags().StringVar(&url, "url", "", "A URL to fetch the aggregate alert file from.")
	cmd.MarkFlagsMutuallyExclusive("file", "url")
	return &cmd
}
