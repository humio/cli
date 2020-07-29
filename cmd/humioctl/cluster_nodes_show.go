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

	"github.com/humio/cli/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func newClusterNodesShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show the information about the a Humio node [Root Only]",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			nodeID := args[0]
			id, parseErr := strconv.Atoi(nodeID)
			exitOnError(cmd, parseErr, "could not parse node id")

			client := NewApiClient(cmd)
			node, apiErr := client.ClusterNodes().Get(id)
			exitOnError(cmd, apiErr, "error fetching node information")
			printClusterNodeInfo(cmd, node)
			cmd.Println()
		},
	}

	return cmd
}

func printClusterNodeInfo(cmd *cobra.Command, node api.ClusterNode) {
	data := [][]string{
		[]string{"ID", strconv.Itoa(node.Id)},
		[]string{"Name", node.Name},
		[]string{"URI", node.Uri},
		[]string{"UUID", node.Uuid},
		[]string{"Cluster info age (Seconds)", fmt.Sprintf("%.3f", node.ClusterInfoAgeSeconds)},
		[]string{"Inbound segment (Size)", ByteCountDecimal(int64(node.InboundSegmentSize))},
		[]string{"Outbound segment (Size)", ByteCountDecimal(int64(node.OutboundSegmentSize))},
		[]string{"Storage Divergence (Size)", ByteCountDecimal(int64(node.StorageDivergence))},
		[]string{"Can be safely unregistered", strconv.FormatBool(node.CanBeSafelyUnregistered)},
		[]string{"Current size", ByteCountDecimal(int64(node.CurrentSize))},
		[]string{"Primary size", ByteCountDecimal(int64(node.PrimarySize))},
		[]string{"Secondary size", ByteCountDecimal(int64(node.SecondarySize))},
		[]string{"Total size of primary", ByteCountDecimal(int64(node.TotalSizeOfPrimary))},
		[]string{"Total size of secondary", ByteCountDecimal(int64(node.TotalSizeOfSecondary))},
		[]string{"Free on primary", ByteCountDecimal(int64(node.FreeOnPrimary))},
		[]string{"Free on secondary", ByteCountDecimal(int64(node.FreeOnSecondary))},
		[]string{"WIP size", ByteCountDecimal(int64(node.WipSize))},
		[]string{"Target size", ByteCountDecimal(int64(node.TargetSize))},
		[]string{"Reapply target size", ByteCountDecimal(int64(node.Reapply_targetSize))},
		[]string{"Solitary segment size", ByteCountDecimal(int64(node.SolitarySegmentSize))},
		[]string{"Is available", strconv.FormatBool(node.IsAvailable)},
		[]string{"Last heartbeat", node.LastHeartbeat},
	}

	w := tablewriter.NewWriter(cmd.OutOrStdout())
	w.AppendBulk(data)
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	w.Render()
}
