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

package cmd

import (
	"fmt"

	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

func newIngestTokensListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [flags] <repo>",
		Short: "List all ingest tokens in a repository.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			repo := args[0]

			// Get the HTTP client
			client := NewApiClient(cmd)

			tokens, err := client.IngestTokens().List(repo)

			if err != nil {
				return fmt.Errorf("Error fetching token list: %s", err)
			}

			var output []string
			output = append(output, "Name | Token | Assigned Parser")
			for i := 0; i < len(tokens); i++ {
				token := tokens[i]
				output = append(output, fmt.Sprintf("%v | %v | %v", token.Name, token.Token, valueOrEmpty(token.AssignedParser)))
			}

			table := columnize.SimpleFormat(output)

			fmt.Println()
			fmt.Println(table)
			fmt.Println()

			return nil
		},
	}

	return cmd
}
