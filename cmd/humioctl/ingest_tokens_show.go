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
	"github.com/humio/cli/internal/format"
	"github.com/spf13/cobra"
)

func newIngestTokensShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show [flags] <repo> <token-name>",
		Short: "Show details about an ingest-token in a repository.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repo := args[0]
			name := args[1]
			client := NewApiClient(cmd)

			ingestToken, err := client.IngestTokens().Get(repo, name)
			exitOnError(cmd, err, "Error fetching ingest token")

			details := [][]format.Value{
				{format.String("Name"), format.String(ingestToken.Name)},
				{format.String("Token"), format.String(ingestToken.Token)},
				{format.String("Assigned Parser"), format.String(ingestToken.AssignedParser)},
			}
			printDetailsTable(cmd, details)
		},
	}

	return &cmd
}
