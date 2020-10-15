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

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"sort"
)

func newClusterNodesListCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "list [flags]",
		Short: "List cluster nodes [Root Only]",
		Args:  cobra.ExactArgs(0),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			client := NewApiClient(cmd)

			nodes, apiErr := client.ClusterNodes().List()
			if apiErr != nil {
				return nil, fmt.Errorf("error fetching cluster nodes: %w", apiErr)
			}

			sort.Slice(nodes, func(i, j int) bool {
				return nodes[i].Name < nodes[j].Name
			})

			return nodes, nil
		}),
	}

	return &cmd
}
