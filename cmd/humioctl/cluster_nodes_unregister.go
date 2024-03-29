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
	"strconv"

	"github.com/spf13/cobra"
)

func newClusterNodesUnregisterCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "unregister [flags] <nodeID>",
		Short: "Unregister (remove) a node from the cluster [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			node, err := strconv.Atoi(args[0])
			exitOnError(cmd, err, "Not valid node id")

			err = client.ClusterNodes().Unregister(node, false)
			exitOnError(cmd, err, "Error removing cluster node")

			fmt.Fprintf(cmd.OutOrStdout(), "Successfully unregistered node %q\n", node)
		},
	}

	return &cmd
}
