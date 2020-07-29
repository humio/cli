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

	"github.com/humio/cli/api"
	"github.com/olekukonko/tablewriter"
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

	return cmd
}

func printRepoTable(cmd *cobra.Command, repo api.Repository) {

	data := [][]string{
		[]string{"Name", repo.Name},
		[]string{"Description", repo.Description},
		[]string{"Space Used", ByteCountDecimal(repo.SpaceUsed)},
		[]string{"Ingest Retention (Size)", ByteCountDecimal(int64(repo.IngestRetentionSizeGB * 1e9))},
		[]string{"Storage Retention (Size)", ByteCountDecimal(int64(repo.StorageRetentionSizeGB * 1e9))},
		[]string{"Retention (Days)", fmt.Sprintf("%d", int64(repo.RetentionDays))},
	}

	w := tablewriter.NewWriter(cmd.OutOrStdout())
	w.AppendBulk(data)
	w.SetBorder(false)
	w.SetColumnSeparator(":")
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	w.Render()
}
