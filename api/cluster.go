package api

import "github.com/shurcooL/graphql"

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
	StorageDivergence       float64
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
	Reapply_targetSize      float64
	SolitarySegmentSize     float64
	IsAvailable             bool
	LastHeartbeat           string
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
	StoragePartitions                   []StoragePartition
}

func (c *Client) Clusters() *Clusters { return &Clusters{client: c} }

func (c *Clusters) Get() (Cluster, error) {
	var q struct {
		Cluster Cluster
	}

	graphqlErr := c.client.Query(&q, nil)

	return q.Cluster, graphqlErr
}

type StoragePartitionInput struct {
	ID      graphql.Int   `json:"id"`
	NodeIDs []graphql.Int `json:"nodeIds"`
}

type IngestPartitionInput struct {
	ID      graphql.Int   `json:"id"`
	NodeIDs []graphql.Int `json:"nodeIds"`
}

func (c *Clusters) UpdateStoragePartitionScheme(desired []StoragePartitionInput) error {
	var m struct {
		UpdateStoragePartitionScheme struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateStoragePartitionScheme(partitions: $partitions)"`
	}

	variables := map[string]interface{}{
		"partitions": desired,
	}

	graphqlErr := c.client.Mutate(&m, variables)

	return graphqlErr
}

func (c *Clusters) UpdateIngestPartitionScheme(desired []IngestPartitionInput) error {
	var m struct {
		UpdateStoragePartitionScheme struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateIngestPartitionScheme(partitions: $partitions)"`
	}

	variables := map[string]interface{}{
		"partitions": desired,
	}

	graphqlErr := c.client.Mutate(&m, variables)

	return graphqlErr
}

func (c *Clusters) StartDataRedistribution() error {
	var m struct {
		StartDataRedistribution struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"startDataRedistribution"`
	}

	graphqlErr := c.client.Mutate(&m, nil)

	return graphqlErr
}

func (c *Clusters) ClusterMoveStorageRouteAwayFromNode(nodeID int) error {
	var m struct {
		ClusterMoveStorageRouteAwayFromNode struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"clusterMoveStorageRouteAwayFromNode(nodeID: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.Int(nodeID),
	}

	graphqlErr := c.client.Mutate(&m, variables)

	return graphqlErr
}

func (c *Clusters) ClusterMoveIngestRoutesAwayFromNode(nodeID int) error {
	var m struct {
		ClusterMoveIngestRoutesAwayFromNode struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"clusterMoveIngestRoutesAwayFromNode(nodeID: $id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.Int(nodeID),
	}

	graphqlErr := c.client.Mutate(&m, variables)

	return graphqlErr
}
