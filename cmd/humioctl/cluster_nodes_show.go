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
		{"ID", strconv.Itoa(node.Id)},
		{"Name", node.Name},
		{"URI", node.Uri},
		{"UUID", node.Uuid},
		{"Cluster info age (Seconds)", fmt.Sprintf("%.3f", node.ClusterInfoAgeSeconds)},
		{"Inbound segment (Size)", ByteCountDecimal(int64(node.InboundSegmentSize))},
		{"Outbound segment (Size)", ByteCountDecimal(int64(node.OutboundSegmentSize))},
		{"Storage Divergence (Size)", ByteCountDecimal(int64(node.StorageDivergence))},
		{"Can be safely unregistered", strconv.FormatBool(node.CanBeSafelyUnregistered)},
		{"Current size", ByteCountDecimal(int64(node.CurrentSize))},
		{"Primary size", ByteCountDecimal(int64(node.PrimarySize))},
		{"Secondary size", ByteCountDecimal(int64(node.SecondarySize))},
		{"Total size of primary", ByteCountDecimal(int64(node.TotalSizeOfPrimary))},
		{"Total size of secondary", ByteCountDecimal(int64(node.TotalSizeOfSecondary))},
		{"Free on primary", ByteCountDecimal(int64(node.FreeOnPrimary))},
		{"Free on secondary", ByteCountDecimal(int64(node.FreeOnSecondary))},
		{"WIP size", ByteCountDecimal(int64(node.WipSize))},
		{"Target size", ByteCountDecimal(int64(node.TargetSize))},
		{"Reapply target size", ByteCountDecimal(int64(node.Reapply_targetSize))},
		{"Solitary segment size", ByteCountDecimal(int64(node.SolitarySegmentSize))},
		{"Is available", strconv.FormatBool(node.IsAvailable)},
		{"Last heartbeat", node.LastHeartbeat},
	}

	w := tablewriter.NewWriter(cmd.OutOrStdout())
	w.AppendBulk(data)
	w.SetColumnAlignment([]int{tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	w.Render()
}
