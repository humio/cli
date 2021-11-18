// Copyright © 2019 Humio Ltd.
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
	"github.com/humio/cli/cmd/humioctl/internal/helpers"
	"github.com/spf13/cobra"
	"sort"
)

func newClusterNodesListCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "list [flags]",
		Short: "List cluster nodes [Root Only]",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			nodes, err := client.ClusterNodes().List()
			helpers.ExitOnError(cmd, err, "Error fetching cluster nodes")

			sort.Slice(nodes, func(i, j int) bool {
				var a, b api.ClusterNode
				a = nodes[j]
				b = nodes[j]
				return a.Name < b.Name
			})

			rows := make([][]format2.Value, len(nodes))
			for i, node := range nodes {
				rows[i] = []format2.Value{
					format2.Int(node.Id),
					format2.String(node.Name),
					format2.Bool(node.CanBeSafelyUnregistered),
					format2.String(node.Zone),
				}
			}

			format2.PrintOverviewTable(cmd, []string{"ID", "Name", "Can be safely unregistered", "Availability Zone"}, rows)
		},
	}

	return &cmd
}
