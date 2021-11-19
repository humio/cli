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

package humioctl

import (
	"fmt"
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func newNotifiersInstallCmd() *cobra.Command {
	var force bool
	var filePath, url, name string

	cmd := cobra.Command{
		Use:   "install [flags] <repo-or-view>",
		Short: "Installs a notifier in a view",
		Long: `Install a notifier from a URL or from a local file.

The install command allows you to install notifiers from a URL or from a local file, e.g.

  $ humioctl notifiers install viewName --name notifierName --url=https://example.com/acme/notifier.yaml

  $ humioctl notifiers install viewName --name notifierName --file=./notifier.yaml

  $ humioctl notifiers install viewName --file=./notifier.yaml

By default 'install' will not override existing parsers with the same name.
Use the --force flag to update existing parsers with conflicting names.
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var content []byte
			var err error

			// Check that we got the right number of argument
			// if we only got <view> you must supply --file or --url.
			if l := len(args); l == 1 {
				if filePath != "" {
					content, err = getNotifierFromFile(filePath)
				} else if url != "" {
					content, err = getURLNotifier(url)
				} else {
					cmd.PrintErrln("You must specify a path using --file or --url")
					os.Exit(1)
				}
			}
			helpers.ExitOnError(cmd, err, "Failed to load the notifier")

			client := NewApiClient(cmd)
			viewName := args[0]

			notifier := api.Notifier{}
			err = yaml.Unmarshal(content, &notifier)
			helpers.ExitOnError(cmd, err, "The notifier's format was invalid")

			if name != "" {
				notifier.Name = name
			}

			_, err = client.Notifiers().Add(viewName, &notifier, force)
			helpers.ExitOnError(cmd, err, "Error installing notifier")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully installed notifier with name: %q\n", notifier.Name)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overrides any notifier with the same name. This can be used for updating notifier that are already installed. (See --name)")
	cmd.Flags().StringVar(&filePath, "file", "", "The local file path to the notifier to install.")
	cmd.Flags().StringVar(&url, "url", "", "A URL to fetch the notifier file from.")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Install the notifier under a specific name, ignoring the `name` attribute in the notifier file.")

	return &cmd
}

func getNotifierFromFile(filePath string) ([]byte, error) {
	// #nosec G304
	return ioutil.ReadFile(filePath)
}

func getURLNotifier(url string) ([]byte, error) {
	// #nosec G107
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
