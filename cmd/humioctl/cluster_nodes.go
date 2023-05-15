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
	"github.com/humio/cli/api"
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
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
	details := [][]format.Value{
		{format.String("ID"), format.Int(node.Id)},
		{format.String("Name"), format.String(node.Name)},
		{format.String("URI"), format.String(node.Uri)},
		{format.String("UUID"), format.String(node.Uuid)},
		{format.String("Cluster info age (Seconds)"), format.Float(node.ClusterInfoAgeSeconds)},
		{format.String("Inbound segment (Size)"), ByteCountDecimal(node.InboundSegmentSize)},
		{format.String("Outbound segment (Size)"), ByteCountDecimal(node.OutboundSegmentSize)},
		{format.String("Can be safely unregistered"), format.Bool(node.CanBeSafelyUnregistered)},
		{format.String("Current size"), ByteCountDecimal(node.CurrentSize)},
		{format.String("Primary size"), ByteCountDecimal(node.PrimarySize)},
		{format.String("Secondary size"), ByteCountDecimal(node.SecondarySize)},
		{format.String("Total size of primary"), ByteCountDecimal(node.TotalSizeOfPrimary)},
		{format.String("Total size of secondary"), ByteCountDecimal(node.TotalSizeOfSecondary)},
		{format.String("Free on primary"), ByteCountDecimal(node.FreeOnPrimary)},
		{format.String("Free on secondary"), ByteCountDecimal(node.FreeOnSecondary)},
		{format.String("WIP size"), ByteCountDecimal(node.WipSize)},
		{format.String("Target size"), ByteCountDecimal(node.TargetSize)},
		{format.String("Solitary segment size"), ByteCountDecimal(node.SolitarySegmentSize)},
		{format.String("Is available"), format.Bool(node.IsAvailable)},
		{format.String("Last heartbeat"), format.String(node.LastHeartbeat)},
		{format.String("Availability Zone"), format.String(node.Zone)},
	}

	printDetailsTable(cmd, details)
}
