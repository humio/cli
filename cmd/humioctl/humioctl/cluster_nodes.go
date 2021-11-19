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
	details := [][]format2.Value{
		{format2.String("ID"), format2.Int(node.Id)},
		{format2.String("Name"), format2.String(node.Name)},
		{format2.String("URI"), format2.String(node.Uri)},
		{format2.String("UUID"), format2.String(node.Uuid)},
		{format2.String("Cluster info age (Seconds)"), format2.Float(node.ClusterInfoAgeSeconds)},
		{format2.String("Inbound segment (Size)"), format2.ByteCountDecimal(node.InboundSegmentSize)},
		{format2.String("Outbound segment (Size)"), format2.ByteCountDecimal(node.OutboundSegmentSize)},
		{format2.String("Storage Divergence (Size)"), format2.ByteCountDecimal(node.StorageDivergence)},
		{format2.String("Can be safely unregistered"), format2.Bool(node.CanBeSafelyUnregistered)},
		{format2.String("Current size"), format2.ByteCountDecimal(node.CurrentSize)},
		{format2.String("Primary size"), format2.ByteCountDecimal(node.PrimarySize)},
		{format2.String("Secondary size"), format2.ByteCountDecimal(node.SecondarySize)},
		{format2.String("Total size of primary"), format2.ByteCountDecimal(node.TotalSizeOfPrimary)},
		{format2.String("Total size of secondary"), format2.ByteCountDecimal(node.TotalSizeOfSecondary)},
		{format2.String("Free on primary"), format2.ByteCountDecimal(node.FreeOnPrimary)},
		{format2.String("Free on secondary"), format2.ByteCountDecimal(node.FreeOnSecondary)},
		{format2.String("WIP size"), format2.ByteCountDecimal(node.WipSize)},
		{format2.String("Target size"), format2.ByteCountDecimal(node.TargetSize)},
		{format2.String("Reapply target size"), format2.ByteCountDecimal(node.Reapply_targetSize)},
		{format2.String("Solitary segment size"), format2.ByteCountDecimal(node.SolitarySegmentSize)},
		{format2.String("Is available"), format2.Bool(node.IsAvailable)},
		{format2.String("Last heartbeat"), format2.String(node.LastHeartbeat)},
		{format2.String("Availability Zone"), format2.String(node.Zone)},
	}

	format2.PrintDetailsTable(cmd, details)
}
