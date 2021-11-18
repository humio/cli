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

package humioctl

import (
	"github.com/humio/cli/api"
	format2 "github.com/humio/cli/cmd/humioctl/internal/format"
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
	details := [][]format2.Value{
		{format2.String("ID"), format2.String(repo.ID)},
		{format2.String("Name"), format2.String(repo.Name)},
		{format2.String("Description"), format2.String(repo.Description)},
		{format2.String("Space Used"), format2.ByteCountDecimal(repo.SpaceUsed)},
		{format2.String("Ingest Retention (Size)"), format2.ByteCountDecimal(repo.IngestRetentionSizeGB * 1e9)},
		{format2.String("Storage Retention (Size)"), format2.ByteCountDecimal(repo.StorageRetentionSizeGB * 1e9)},
		{format2.String("Retention (Days)"), format2.Int(repo.RetentionDays)},
	}

	format2.PrintDetailsTable(cmd, details)
}
