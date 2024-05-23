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
	var allowDataDeletionFlag, enableS3ArchivingFlag, disableS3ArchivingFlag bool
	var descriptionFlag, s3ArchivingBucketFlag, s3ArchivingRegionFlag, s3ArchivingFormatFlag stringPtrFlag
	var retentionTimeFlag, ingestSizeBasedRetentionFlag, storageSizeBasedRetentionFlag float64PtrFlag

	cmd := cobra.Command{
		Use:   "update [flags] <repo>",
		Short: "Updates the settings of a repository",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			repoName := args[0]
			client := NewApiClient(cmd)

			if descriptionFlag.value == nil && retentionTimeFlag.value == nil && ingestSizeBasedRetentionFlag.value == nil && storageSizeBasedRetentionFlag.value == nil &&
				s3ArchivingBucketFlag.value == nil && s3ArchivingRegionFlag.value == nil && s3ArchivingFormatFlag.value == nil {
				exitOnError(cmd, fmt.Errorf("you must specify at least one flag to update"), "Nothing specified to update")
			}
			if descriptionFlag.value != nil {
				err := client.Repositories().UpdateDescription(repoName, *descriptionFlag.value)
				exitOnError(cmd, err, "Error updating repository description")
			}
			if retentionTimeFlag.value != nil {
				err := client.Repositories().UpdateTimeBasedRetention(repoName, *retentionTimeFlag.value, allowDataDeletionFlag)
				exitOnError(cmd, err, "Error updating repository retention time in days")
			}
			if ingestSizeBasedRetentionFlag.value != nil {
				err := client.Repositories().UpdateIngestBasedRetention(repoName, *ingestSizeBasedRetentionFlag.value, allowDataDeletionFlag)
				exitOnError(cmd, err, "Error updating repository ingest size based retention")
			}
			if storageSizeBasedRetentionFlag.value != nil {
				err := client.Repositories().UpdateStorageBasedRetention(repoName, *storageSizeBasedRetentionFlag.value, allowDataDeletionFlag)
				exitOnError(cmd, err, "Error updating repository storage size based retention")
			}

			if s3ArchivingBucketFlag.value != nil && s3ArchivingRegionFlag.value != nil && s3ArchivingFormatFlag.value != nil {
				err := client.Repositories().UpdateS3ArchivingConfiguration(repoName, *s3ArchivingBucketFlag.value, *s3ArchivingRegionFlag.value, *s3ArchivingFormatFlag.value)
				exitOnError(cmd, err, "Error updating S3 archiving configuration")
			}

			if disableS3ArchivingFlag == true {
				err := client.Repositories().DisableS3Archiving(repoName)
				exitOnError(cmd, err, "Error disabling S3 archiving")
			}

			if enableS3ArchivingFlag == true {
				err := client.Repositories().EnableS3Archiving(repoName)
				exitOnError(cmd, err, "Error enabling S3 archiving")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully updated repository %q\n", repoName)
		},
	}

	cmd.Flags().BoolVar(&allowDataDeletionFlag, "allow-data-deletion", false, "Allow changing retention settings for a non-empty repository")
	cmd.Flags().Var(&descriptionFlag, "description", "The description of the repository.")
	cmd.Flags().Var(&retentionTimeFlag, "retention-time", "The retention time in days for the repository.")
	cmd.Flags().Var(&ingestSizeBasedRetentionFlag, "ingest-size-based-retention", "The ingest size based retention for the repository.")
	cmd.Flags().Var(&storageSizeBasedRetentionFlag, "storage-size-based-retention", "The storage size based retention for the repository.")
	cmd.Flags().BoolVar(&enableS3ArchivingFlag, "enable-s3-archiving", false, "Enable S3 Archiving")
	cmd.Flags().BoolVar(&disableS3ArchivingFlag, "disable-s3-archiving", false, "Disable S3 Archiving")
	cmd.Flags().Var(&s3ArchivingBucketFlag, "s3-archiving-bucket", "The name of the bucket to be used for S3 Archiving")
	cmd.Flags().Var(&s3ArchivingRegionFlag, "s3-archiving-region", "The S3 region to be used for S3 Archiving")
	cmd.Flags().Var(&s3ArchivingFormatFlag, "s3-archiving-format", "The S3 archiving format to be used for S3 Archiving. Formats: RAW, NDJSON")
	cmd.MarkFlagsRequiredTogether("s3-archiving-bucket", "s3-archiving-region", "s3-archiving-format")
	cmd.MarkFlagsMutuallyExclusive("enable-s3-archiving", "disable-s3-archiving")

	return &cmd
}
