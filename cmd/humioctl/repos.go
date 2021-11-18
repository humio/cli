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
	"github.com/humio/cli/cmd/internal/format"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
)

func newReposCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "Manage repositories",
	}

	cmd.AddCommand(newReposShowCmd())
	cmd.AddCommand(newReposListCmd())
	cmd.AddCommand(newReposCreateCmd())
	cmd.AddCommand(newReposUpdateCmd())
	cmd.AddCommand(newReposDeleteCmd())
	cmd.AddCommand(newReposUpdateUserGroupCmd())

	return cmd
}

func printRepoDetailsTable(cmd *cobra.Command, repo api.Repository) {
	details := [][]format.Value{
		{format.String("ID"), format.String(repo.ID)},
		{format.String("Name"), format.String(repo.Name)},
		{format.String("Description"), format.String(repo.Description)},
		{format.String("Space Used"), format.ByteCountDecimal(repo.SpaceUsed)},
		{format.String("Ingest Retention (Size)"), format.ByteCountDecimal(repo.IngestRetentionSizeGB * 1e9)},
		{format.String("Storage Retention (Size)"), format.ByteCountDecimal(repo.StorageRetentionSizeGB * 1e9)},
		{format.String("Retention (Days)"), format.Int(repo.RetentionDays)},
	}

	format.PrintDetailsTable(cmd, details)
}
