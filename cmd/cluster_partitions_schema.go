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

package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/humio/cli/api"
	"github.com/spf13/cobra"
)

func newClusterPartitionsSchemaCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "list [flags]",
		Short: "List cluster nodes [Root Only]",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			client := NewApiClient(cmd)

			nodes, apiErr := client.ClusterNodes().List()
			exitOnError(cmd, apiErr, "error fetching cluster nodes")

			var storagePartitions []api.StoragePartition
			var spErr error
			storagePartitions, spErr = client.ClusterPartitions().GetStoragePartitions()
			if spErr != nil {
				fmt.Println(fmt.Sprintf("%v", &storagePartitions))

			}
			// fmt.Println(fmt.Sprintf("%s", &storagePartitions))
			// var jsonString []byte
			// jsonString, _ = json.Marshal(&storagePartitions)
			// fmt.Println(fmt.Sprintf("%s", jsonString))

			sort.Slice(nodes, func(i, j int) bool {
				mi, mj := nodes[i], nodes[j]
				aSplit := strings.SplitAfter(mi.UUID, "-")
				bSplit := strings.SplitAfter(mj.UUID, "-")
				aRackFull := aSplit[len(aSplit)-1]
				bRackFull := bSplit[len(bSplit)-1]
				aRackSplit := strings.SplitAfter(aRackFull, "_")
				bRackSplit := strings.SplitAfter(bRackFull, "_")
				arackID := aRackSplit[1]
				brackID := bRackSplit[1]
				arack := aRackSplit[0]
				brack := bRackSplit[0]
				// return arackID < brackID
				switch {
				case arackID != brackID:
					return arackID < brackID
				default:
					return arack < brack
				}

			})
			var storageNodeIds []int
			for _, node := range nodes {
				storageNodeIds = append(storageNodeIds, node.Id)
			}

			var ps []api.StoragePartition

			var partitionCount int
			partitionCount = 24
			var replicas int
			replicas = 2

			for p := 0; p < partitionCount; p++ {
				var nodeIds []int
				for r := 0; r < replicas; r++ {
					idx := (p + r) % len(storageNodeIds)
					nodeIds = append(nodeIds, storageNodeIds[idx])
				}
				ps = append(ps, api.StoragePartition{Id: p, NodeIds: nodeIds})
			}
			var jsonString []byte
			jsonString, _ = json.Marshal(&ps)
			fmt.Println(fmt.Sprintf("%s", jsonString))

			// fmt.Println(fmt.Sprintf("%v", &nodes))
			// rows := make([][]string, len(nodes))
			// for i, node := range nodes {
			// 	rows[i] = []string{strconv.Itoa(node.Id), node.Name, node.UUID, strconv.FormatBool(node.CanBeSafelyUnregistered)}
			// }

			// w := tablewriter.NewWriter(cmd.OutOrStdout())
			// w.SetHeader([]string{"ID", "Name", "UUID", "Can be safely unregistered"})
			// w.AppendBulk(rows)
			// w.SetBorder(false)

			// w.Render()
			cmd.Println()
		},
	}

	return &cmd
}
