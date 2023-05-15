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
	"io"
	"net/http"
	"os"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newActionsInstallCmd() *cobra.Command {
	var filePath, url, name string

	cmd := cobra.Command{
		Use:   "install [flags] <repo-or-view>",
		Short: "Installs an action in a view",
		Long: `Install an action from a URL or from a local file.

The install command allows you to install actions from a URL or from a local file, e.g.

  $ humioctl actions install viewName --name actionName --url=https://example.com/acme/action.yaml

  $ humioctl actions install viewName --name actionName --file=./action.yaml

  $ humioctl actions install viewName --file=./action.yaml

By default 'install' will not override existing actions with the same name.
Use the --force flag to update existing actions with conflicting names.
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var content []byte
			var err error

			// Check that we got the right number of argument
			// if we only got <view> you must supply --file or --url.
			if l := len(args); l == 1 {
				if filePath != "" {
					content, err = getActionFromFile(filePath)
				} else if url != "" {
					content, err = getURLAction(url)
				} else {
					cmd.PrintErrln("You must specify a path using --file or --url")
					os.Exit(1)
				}
			}
			exitOnError(cmd, err, "Failed to load the action")

			client := NewApiClient(cmd)
			viewName := args[0]

			action := api.Action{}
			err = yaml.Unmarshal(content, &action)
			exitOnError(cmd, err, "The action's format was invalid")

			if name != "" {
				action.Name = name
			}

			_, err = client.Actions().Add(viewName, &action)
			exitOnError(cmd, err, "Error installing action")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully installed action with name: %q\n", action.Name)
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "The local file path to the action to install.")
	cmd.Flags().StringVar(&url, "url", "", "A URL to fetch the action file from.")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Install the action under a specific name, ignoring the `name` attribute in the action file.")

	return &cmd
}

func getActionFromFile(filePath string) ([]byte, error) {
	// #nosec G304
	return os.ReadFile(filePath)
}

func getURLAction(url string) ([]byte, error) {
	// #nosec G107
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()
	return io.ReadAll(response.Body)
}
