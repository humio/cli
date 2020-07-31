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

	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

func newIngestTokensUpdateCmd() *cobra.Command {
	var parserName string

	cmd := &cobra.Command{
		Use:   "update [flags] <repository-name> <token-name>",
		Short: "Update an ingest token to a repository.",
		Long: `Updates the parser of an ingest token with name <token name> in repository <repo>.

You can associate a parser with the ingest token using the --parser flag.
Assigning a parser will make all data sent to Humio using this ingest token
use the assigned parser at ingest time.

If parser is not specified, the ingest token will not be associated with a parser.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repositoryName := args[0]
			tokenName := args[1]

			// Get the HTTP client
			client := NewApiClient(cmd)

			token, err := client.IngestTokens().Update(repositoryName, tokenName, parserName)

			if err != nil {
				return fmt.Errorf("Error updating ingest token: %s", err)
			}

			var output []string
			output = append(output, "Name | Token | Assigned Parser")
			output = append(output, fmt.Sprintf("%v | %v | %v", token.Name, token.Token, valueOrEmpty(token.AssignedParser)))

			table := columnize.SimpleFormat(output)

			cmd.Println()
			cmd.Println(table)
			cmd.Println()

			return nil
		},
	}

	cmd.Flags().StringVarP(&parserName, "parser", "p", "", "Assigns the a parser to the ingest token.")

	return cmd
}
