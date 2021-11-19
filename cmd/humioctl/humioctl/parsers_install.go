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

package humioctl

import (
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func newParsersInstallCmd() *cobra.Command {
	var force bool
	var filePath, url, name string

	cmd := cobra.Command{
		Use:   "install [flags] <repo>",
		Short: "Installs a parser in a repository",
		Long: `Install a parser from a URL or from a local file.

The install command will install parser from a local file or a URL, e.g.

  $ humioctl parsers install --url=https://example.com/acme/parser.yaml

  $ humioctl parsers install --file=./parser.yaml

By default 'install' will not override existing parsers with the same name.
Use the --force flag to update existing parsers with conflicting names.
`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var content []byte
			var err error

			// Check that we got the right number of argument
			// if we only got <repo> you must supply --file or --url.
			if filePath != "" {
				content, err = getParserFromFile(filePath)
			} else if url != "" {
				content, err = getURLParser(url)
			} else {
				cmd.PrintErrf("If you only provide repo you must specify --file or --url\n")
				os.Exit(1)
			}
			helpers.ExitOnError(cmd, err, "Failed to load the parser")

			repositoryName := args[0]
			client := NewApiClient(cmd)

			parser := api.Parser{}
			err = yaml.Unmarshal(content, &parser)
			helpers.ExitOnError(cmd, err, "The parser's format was invalid")

			if name != "" {
				parser.Name = name
			}

			err = client.Parsers().Add(repositoryName, &parser, force)
			helpers.ExitOnError(cmd, err, "Error installing parser")
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overrides any parser with the same name. This can be used for updating parser that are already installed. (See --name)")
	cmd.Flags().StringVar(&filePath, "file", "", "The local file path to the parser to install.")
	cmd.Flags().StringVar(&url, "url", "", "A URL to fetch the parser file from.")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Install the parser under a specific name, ignoring the `name` attribute in the parser file.")

	return &cmd
}

func getParserFromFile(filePath string) ([]byte, error) {
	// #nosec G304
	return ioutil.ReadFile(filePath)
}

func getURLParser(url string) ([]byte, error) {
	// #nosec G107
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
