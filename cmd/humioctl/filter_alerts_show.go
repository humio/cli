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
	"strings"

	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newFilterAlertsShowCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "show <repo-or-view> <name>",
		Short: "Show details about an filter alert in a repository or view.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			repoOrViewName := args[0]
			name := args[1]
			client := NewApiClient(cmd)

			filterAlert, err := client.FilterAlerts().Get(repoOrViewName, name)
			exitOnError(cmd, err, "Error fetching filter alert")

			details := [][]format.Value{
				{format.String("ID"), format.String(filterAlert.ID)},
				{format.String("Name"), format.String(filterAlert.Name)},
				{format.String("Enabled"), format.Bool(filterAlert.Enabled)},
				{format.String("Description"), format.String(filterAlert.Description)},
				{format.String("Query String"), format.String(filterAlert.QueryString)},
				{format.String("Labels"), format.String(strings.Join(filterAlert.Labels, ", "))},
				{format.String("Actions"), format.String(strings.Join(filterAlert.Actions, ", "))},
				{format.String("Last Error"), format.String(filterAlert.LastError)},
				{format.String("Last Error Time"), format.Int(filterAlert.LastErrorTime)},
				{format.String("Last Error"), format.String(filterAlert.LastError)},
				{format.String("Run As User ID"), format.String(filterAlert.RunAsUserID)},
				{format.String("Query Ownership Type"), format.String(filterAlert.QueryOwnershipType)},
			}

			printDetailsTable(cmd, details)
		},
	}

	return &cmd
}
