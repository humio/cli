// Copyright © 2020 Humio Ltd.
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
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

func newAlertsInstallCmd() *cobra.Command {
	var content []byte
	var readErr error
	var force bool
	var filePath, url, name string

	cmd := cobra.Command{
		Use:   "install [flags] <view>",
		Short: "Installs an alert in a view",
		Long: `Install an alert from a URL or from a local file.

The install command allows you to install alerts from a URL or from a local file, e.g.

  $ humioctl alerts install viewName alertName --url=https://example.com/acme/alert.yaml

  $ humioctl alerts install viewName alertName --file=./parser.yaml

  $ humioctl alerts install viewName --file=./alert.yaml

By default 'install' will not override existing alerts with the same name.
Use the --force flag to update existing alerts with conflicting names.
`,
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			// Check that we got the right number of argument
			// if we only got <view> you must supply --file or --url.
			if len(args) == 1 {
				if filePath != "" {
					content, readErr = getAlertFromFile(filePath)
				} else if url != "" {
					content, readErr = getURLAlert(url)
				} else {
					return nil, fmt.Errorf("you must specify a path using --file or --url")
				}
			} else if len(args) != 2 {
				return nil, fmt.Errorf("this command takes one argument: <view>")
			}
			if readErr != nil {
				return nil, fmt.Errorf("failed to load the alert: %w", readErr)
			}

			viewName := args[0]
			alert := api.Alert{}
			alert.Name = name
			yamlErr := yaml.Unmarshal(content, &alert)
			if yamlErr != nil {
				return nil, fmt.Errorf("the alert's format was invalid: %w", yamlErr)
			}

			// Get the HTTP client
			client := NewApiClient(cmd)

			_, installErr := client.Alerts().Add(viewName, &alert, force)
			if installErr != nil {
				return nil, fmt.Errorf("error installing alert: %w", installErr)
			}

			return fmt.Sprintf("Alert %q installed.", alert.Name), nil
		}),
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overrides any alert with the same name. This can be used for updating alert that are already installed. (See --name)")
	cmd.Flags().StringVar(&filePath, "file", "", "The local file path to the alert to install.")
	cmd.Flags().StringVar(&url, "url", "", "A URL to fetch the alert file from.")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Install the alert under a specific name, ignoreing the `name` attribute in the alert file.")

	return &cmd
}

func getAlertFromFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func getURLAlert(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
