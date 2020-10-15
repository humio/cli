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

func newClusterNodesShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show the information about the a Humio node [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: WrapRun(func(cmd *cobra.Command, args []string) (humioResultType, error) {
			nodeID := args[0]
			id, parseErr := strconv.Atoi(nodeID)
			if parseErr != nil {
				return nil, fmt.Errorf("could not parse node id: %w", parseErr)
			}

			client := NewApiClient(cmd)
			node, apiErr := client.ClusterNodes().Get(id)
			if apiErr != nil {
				return nil, fmt.Errorf("error fetching node information: %w", apiErr)
			}

			return node, nil
		}),
	}

	return cmd
}
