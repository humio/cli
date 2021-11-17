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
	"github.com/humio/cli/cmd/internal/format"
	"github.com/spf13/cobra"
)

func newClusterShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show the information about the current Humio cluster",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			cluster, err := client.Clusters().Get()
			exitOnError(cmd, err, "Error fetching cluster information")

			rows := [][]format.Value{
				{
					format.String("Under replicated segment (Size)"),
					ByteCountDecimal(int64(cluster.UnderReplicatedSegmentSize)),
					ByteCountDecimal(int64(cluster.TargetUnderReplicatedSegmentSize))},
				{
					format.String("Over replicated segment (Size)"),
					ByteCountDecimal(int64(cluster.OverReplicatedSegmentSize)),
					ByteCountDecimal(int64(cluster.TargetOverReplicatedSegmentSize))},
				{
					format.String("Missing segment (Size)"),
					ByteCountDecimal(int64(cluster.MissingSegmentSize)),
					ByteCountDecimal(int64(cluster.TargetMissingSegmentSize))},
				{
					format.String("Properly replicated segment (Size)"),
					ByteCountDecimal(int64(cluster.ProperlyReplicatedSegmentSize)),
					ByteCountDecimal(int64(cluster.TargetProperlyReplicatedSegmentSize))},
			}

			printOverviewTable(cmd, []string{"Description", "Current", "Target"}, rows)
		},
	}

	return cmd
}
