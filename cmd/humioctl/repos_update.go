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

	"github.com/spf13/cobra"
)

func newReposUpdateCmd() *cobra.Command {
	var allowDataDeletionFlag bool
	var descriptionFlag stringPtrFlag
	var retentionTimeFlag, ingestSizeBasedRetentionFlag, storageSizeBasedretentionFlag float64PtrFlag

	cmd := cobra.Command{
		Use:   "update",
		Short: "Updates the settings of a repository",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]

			if descriptionFlag.value == nil && retentionTimeFlag.value == nil && ingestSizeBasedRetentionFlag.value == nil && storageSizeBasedretentionFlag.value == nil {
				exitOnError(cmd, fmt.Errorf("you must specify at least one flag to update"), "nothing specifed to update")
			}

			client := NewApiClient(cmd)
			if descriptionFlag.value != nil {
				err := client.Repositories().UpdateDescription(repoName, *descriptionFlag.value)
				exitOnError(cmd, err, "error updating repository description")
			}
			if retentionTimeFlag.value != nil {
				err := client.Repositories().UpdateTimeBasedRetention(repoName, *retentionTimeFlag.value, allowDataDeletionFlag)
				exitOnError(cmd, err, "error updating repository retention time in days")
			}
			if ingestSizeBasedRetentionFlag.value != nil {
				err := client.Repositories().UpdateIngestBasedRetention(repoName, *ingestSizeBasedRetentionFlag.value, allowDataDeletionFlag)
				exitOnError(cmd, err, "error updating repository ingest size based retention")
			}
			if storageSizeBasedretentionFlag.value != nil {
				err := client.Repositories().UpdateStorageBasedRetention(repoName, *storageSizeBasedretentionFlag.value, allowDataDeletionFlag)
				exitOnError(cmd, err, "error updating repository storage size based retention")
			}

			repo, apiErr := client.Repositories().Get(repoName)
			exitOnError(cmd, apiErr, "error fetching repository")
			printRepoTable(cmd, repo)
			fmt.Println()
		},
	}

	cmd.Flags().BoolVar(&allowDataDeletionFlag, "allow-data-deletion", false, "Allow changing retention settings for a non-empty repository")
	cmd.Flags().Var(&descriptionFlag, "description", "The description of the repository.")
	cmd.Flags().Var(&retentionTimeFlag, "retention-time", "The retention time in days for the repository.")
	cmd.Flags().Var(&ingestSizeBasedRetentionFlag, "ingest-size-based-retention", "The ingest size based retention for the repository.")
	cmd.Flags().Var(&storageSizeBasedretentionFlag, "storage-size-based-retention", "The storage size based retention for the repository.")

	return &cmd
}
