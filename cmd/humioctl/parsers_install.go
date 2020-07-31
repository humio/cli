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

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func newParsersInstallCmd() *cobra.Command {
	var content []byte
	var readErr error
	var force bool
	var filePath, url, name string

	cmd := cobra.Command{
		Use:   "install [flags] <repo> <parser>",
		Short: "Installs a parser in a repository",
		Long: `Install a parser from a Humio's community github repository, a URL or
from a local file.

To find all available parsers visit: https://github.com/humio/community/parsers

For instance if you wanted to install an AccessLog parser you could use.

  $ humioctl parsers install accesslog

This would install the parser at: humio/comminity/parsers/accesslog/default.yaml
Since log formats can vary slightly you can install one of the other variations:

  $ humioctl parsers install accesslog/utc

Which would install the humio/community/parsers/accesslog/utc.yaml parser.

The install command will pull parser from GitHub by default. But you can also
install from a local file or a URL, e.g.

  $ humioctl parsers install --url=https://example.com/acme/parser.yaml

  $ humioctl parsers install --file=./parser.yaml

By default 'install' will not override existing parsers with the same name.
Use the --force flag to update existing parsers with conflicting names.
`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check that we got the right number of argument
			// if we only got <repo> you must supply --file or --url.
			if l := len(args); l == 1 {
				if filePath != "" {
					content, readErr = getParserFromFile(filePath)
				} else if url != "" {
					content, readErr = getURLParser(url)
				} else {
					cmd.Println(fmt.Errorf("if you only provide repo you must specify --file or --url"))
					os.Exit(1)
				}
			} else if l := len(args); l != 2 {
				cmd.Println(fmt.Errorf("This command takes one or two arguments: <repo> [parser]"))
				os.Exit(1)
			} else {
				parserName := args[1]
				content, readErr = getGithubParser(parserName)
			}

			exitOnError(cmd, readErr, "Failed to load the parser")

			parser := api.Parser{}
			yamlErr := yaml.Unmarshal(content, &parser)
			exitOnError(cmd, yamlErr, "The parser's format was invalid")

			if name != "" {
				parser.Name = name
			}

			// Get the HTTP client
			client := NewApiClient(cmd)

			reposistoryName := args[0]

			installErr := client.Parsers().Add(reposistoryName, &parser, force)
			exitOnError(cmd, installErr, "error installing parser")
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overrides any parser with the same name. This can be used for updating parser that are already installed. (See --name)")
	cmd.Flags().StringVar(&filePath, "file", "", "The local file path to the parser to install.")
	cmd.Flags().StringVar(&url, "url", "", "A URL to fetch the parser file from.")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Install the parser under a specific name, ignoreing the `name` attribute in the parser file.")

	return &cmd
}

func getParserFromFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func getGithubParser(parserName string) ([]byte, error) {
	url := "https://raw.githubusercontent.com/humio/community/master/parsers/" + parserName + ".yaml"
	return getURLParser(url)
}

func getURLParser(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
