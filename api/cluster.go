package api

import (
	"fmt"
	graphql "github.com/cli/shurcooL-graphql"
	"math"
)

type Clusters struct {
	client *Client
}

type ClusterNode struct {
	Id                      int
	Name                    string
	Uri                     string
	Uuid                    string
	ClusterInfoAgeSeconds   float64
	InboundSegmentSize      float64
	OutboundSegmentSize     float64
	CanBeSafelyUnregistered bool
	CurrentSize             float64
	PrimarySize             float64
	SecondarySize           float64
	TotalSizeOfPrimary      float64
	TotalSizeOfSecondary    float64
	FreeOnPrimary           float64
	FreeOnSecondary         float64
	WipSize                 float64
	TargetSize              float64
	SolitarySegmentSize     float64
	IsAvailable             bool
	LastHeartbeat           string
	// Zone holds the availability zone as configured in the `ZONE` configuration of the Humio server.
	Zone string
}

type IngestPartition struct {
	Id      int
	NodeIds []int
}

type StoragePartition struct {
	Id      int
	NodeIds []int
}

type Cluster struct {
	Nodes                               []ClusterNode
	ClusterInfoAgeSeconds               float64
	UnderReplicatedSegmentSize          float64
	OverReplicatedSegmentSize           float64
	MissingSegmentSize                  float64
	ProperlyReplicatedSegmentSize       float64
	TargetUnderReplicatedSegmentSize    float64
	TargetOverReplicatedSegmentSize     float64
	TargetMissingSegmentSize            float64
	TargetProperlyReplicatedSegmentSize float64
	IngestPartitions                    []IngestPartition
	StoragePartitions                   []StoragePartition // Deprecated: returns dummy data as of LogScale 1.88
}

func (c *Client) Clusters() *Clusters { return &Clusters{client: c} }

func (c *Clusters) Get() (Cluster, error) {
	var query struct {
		Cluster Cluster
	}

	err := c.client.Query(&query, nil)
	return query.Cluster, err
}

type StoragePartitionInput struct {
	ID      graphql.Int   `json:"id"`
	NodeIDs []graphql.Int `json:"nodeIds"`
}

// Deprecated: returns dummy data as of LogScale 1.88
func (c *Clusters) UpdateStoragePartitionScheme(desired []StoragePartitionInput) error {
	var mutation struct {
		UpdateStoragePartitionScheme struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateStoragePartitionScheme(partitions: $partitions)"`
	}

	variables := map[string]interface{}{
		"partitions": desired,
	}

	return c.client.Mutate(&mutation, variables)
}

// Deprecated: returns dummy data as of LogScale 1.80
func (c *Clusters) UpdateIngestPartitionScheme(desired []IngestPartitionInput) error {
	var mutation struct {
		UpdateStoragePartitionScheme struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateIngestPartitionScheme(partitions: $partitions)"`
	}

	variables := map[string]interface{}{
		"partitions": desired,
	}

	return c.client.Mutate(&mutation, variables)
}

// Deprecated: returns dummy data as of LogScale 1.88
func (c *Clusters) StartDataRedistribution() error {
	var mutation struct {
		StartDataRedistribution struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"startDataRedistribution"`
	}

	return c.client.Mutate(&mutation, nil)
}

// Deprecated: returns dummy data as of LogScale 1.88
func (c *Clusters) ClusterMoveStorageRouteAwayFromNode(nodeID int) error {
	var mutation struct {
		ClusterMoveStorageRouteAwayFromNode struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"clusterMoveStorageRouteAwayFromNode(nodeID: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.Int(nodeID),
	}

	return c.client.Mutate(&mutation, variables)
}

// Deprecated: returns dummy data as of LogScale 1.80
func (c *Clusters) ClusterMoveIngestRoutesAwayFromNode(nodeID int) error {
	var mutation struct {
		ClusterMoveIngestRoutesAwayFromNode struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"clusterMoveIngestRoutesAwayFromNode(nodeID: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.Int(nodeID),
	}

	return c.client.Mutate(&mutation, variables)
}

type ClusterNodes struct {
	client *Client
}

func (c *Client) ClusterNodes() *ClusterNodes { return &ClusterNodes{client: c} }

func (n *ClusterNodes) List() ([]ClusterNode, error) {
	var query struct {
		Cluster struct {
			Nodes []ClusterNode
		}
	}

	err := n.client.Query(&query, nil)
	return query.Cluster.Nodes, err
}

func (n *ClusterNodes) Get(nodeID int) (ClusterNode, error) {
	var query struct {
		Cluster struct {
			Nodes []ClusterNode
		}
	}

	err := n.client.Query(&query, nil)
	if err != nil {
		return ClusterNode{}, err
	}

	for _, node := range query.Cluster.Nodes {
		if node.Id == nodeID {
			return node, nil
		}
	}

	return ClusterNode{}, fmt.Errorf("node id not found in cluster")
}

func (n *ClusterNodes) Unregister(nodeID int, force bool) error {
	if nodeID > math.MaxInt32 {
		return fmt.Errorf("node id too large")
	}
	var mutation struct {
		ClusterUnregisterNode struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"clusterUnregisterNode(force: $force, nodeID: $id)"`
	}

	variables := map[string]interface{}{
		"id":    graphql.Int(nodeID),
		"force": graphql.Boolean(force),
	}

	return n.client.Mutate(&mutation, variables)
}

// Deprecated: returns dummy data as of LogScale 1.80
func (c *Clusters) SuggestedIngestPartitions() ([]IngestPartitionInput, error) {
	var query struct {
		Cluster struct {
			SuggestedIngestPartitions []IngestPartitionInput `graphql:"suggestedIngestPartitions"`
		} `graphql:"cluster"`
	}

	err := c.client.Query(&query, nil)
	return query.Cluster.SuggestedIngestPartitions, err
}

// Deprecated: returns dummy data as of LogScale 1.88
func (c *Clusters) SuggestedStoragePartitions() ([]StoragePartitionInput, error) {
	var query struct {
		Cluster struct {
			SuggestedStoragePartitions []StoragePartitionInput `graphql:"suggestedStoragePartitions"`
		} `graphql:"cluster"`
	}

	err := c.client.Query(&query, nil)
	return query.Cluster.SuggestedStoragePartitions, err
}
