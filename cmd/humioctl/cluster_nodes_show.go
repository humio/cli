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
	"strconv"

	"github.com/spf13/cobra"
)

func newClusterNodesShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show the information about the a Humio node [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			nodeID := args[0]
			client := NewApiClient(cmd)

			id, err := strconv.Atoi(nodeID)
			exitOnError(cmd, err, "Could not parse node id")

			node, err := client.ClusterNodes().Get(id)
			exitOnError(cmd, err, "Error fetching node information")

			printClusterNodeDetailsTable(cmd, node)
		},
	}

	return cmd
}
