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
	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
	"strconv"
)

func newClusterNodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Manage cluster nodes [Root Only]",
	}

	cmd.AddCommand(newClusterNodesListCmd())
	cmd.AddCommand(newClusterNodesShowCmd())
	cmd.AddCommand(newClusterNodesUnregisterCmd())

	return cmd
}

func printClusterNodeDetailsTable(cmd *cobra.Command, node api.ClusterNode) {
	details := [][]string{
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
		{"Availability Zone", node.Zone},
	}

	printDetailsTable(cmd, details)
}
