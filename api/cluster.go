package api

import (
	"fmt"
	"math"

	graphql "github.com/cli/shurcooL-graphql"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Clusters struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type IngestPartition struct {
	Id      int
	NodeIds []int
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Clusters() *Clusters { return &Clusters{client: c} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Clusters) Get() (Cluster, error) {
	var query struct {
		Cluster Cluster
	}

	err := c.client.Query(&query, nil)
	return query.Cluster, err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type ClusterNodes struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) ClusterNodes() *ClusterNodes { return &ClusterNodes{client: c} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (n *ClusterNodes) List() ([]ClusterNode, error) {
	var query struct {
		Cluster struct {
			Nodes []ClusterNode
		}
	}

	err := n.client.Query(&query, nil)
	return query.Cluster.Nodes, err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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
