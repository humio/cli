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
	"github.com/spf13/cobra"
)

func newIngestTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest-tokens [flags]",
		Short: "Manage ingest tokens",
		Long: `Ingest tokens, unlike the more general API tokens, can only be used for ingestion of data.

You can also assign a parser to an ingest token, allowing you to configure how Humio parses incoming data
without having to change anything on sender/client.`,
	}

	cmd.AddCommand(newIngestTokensAddCmd())
	cmd.AddCommand(newIngestTokensUpdateCmd())
	cmd.AddCommand(newIngestTokensRemoveCmd())
	cmd.AddCommand(newIngestTokensListCmd())
	cmd.AddCommand(newIngestTokensShowCmd())

	return cmd
}
