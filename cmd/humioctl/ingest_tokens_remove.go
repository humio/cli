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

	"github.com/spf13/cobra"
)

func newIngestTokensRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "remove [flags] <repo> <token-name>",
		Short:     "Removes an ingest token.",
		Long:      `Removes the ingest token with name '<token-name>' from the repository with name '<repo>'.`,
		ValidArgs: []string{"repo", "token-name"},
		Args:      cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			tokenName := args[1]
			client := NewApiClient(cmd)

			err := client.IngestTokens().Remove(repoName, tokenName)
			exitOnError(cmd, err, "Error removing ingest token")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully removed ingest token %q from repository %q\n", tokenName, repoName)
		},
	}

	return cmd
}
