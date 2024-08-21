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

func newReposCreateCmd() *cobra.Command {
	var descriptionFlag string
	var retentionTimeFlag int64
	var ingestSizeBasedRetentionFlag, storageSizeBasedRetentionFlag float64

	cmd := cobra.Command{
		Use:   "create <repo>",
		Short: "Create a repository.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			client := NewApiClient(cmd)

			err := client.Repositories().Create(repoName,
				descriptionFlag,
				retentionTimeFlag,
				ingestSizeBasedRetentionFlag,
				storageSizeBasedRetentionFlag,
			)
			exitOnError(cmd, err, "Error creating repository")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully created repo %s\n", repoName)

		},
	}

	cmd.Flags().StringVar(&descriptionFlag, "description", "", "The description of the repository.")
	cmd.Flags().Int64Var(&retentionTimeFlag, "retention-time", 0, "The retention time in days for the repository.")
	cmd.Flags().Float64Var(&ingestSizeBasedRetentionFlag, "retention-ingest", 0, "The ingest size based retention in GB for the repository.")
	cmd.Flags().Float64Var(&storageSizeBasedRetentionFlag, "retention-storage", 0, "The storage size based retention in GB for the repository.")

	return &cmd
}
