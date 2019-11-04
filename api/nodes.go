package api

import (
	"fmt"

	"github.com/shurcooL/graphql"
)

type Nodes struct {
	client *Client
}

func (n *Client) Nodes() *Nodes { return &Nodes{client: n} }

func (n *Nodes) List() ([]ClusterNode, error) {
	var q struct {
		Cluster struct {
			Nodes []ClusterNode
		}
	}

	graphqlErr := n.client.Query(&q, nil)

	return q.Cluster.Nodes, graphqlErr
}

func (n *Nodes) Get(nodeID int) (ClusterNode, error) {
	var q struct {
		Cluster struct {
			Nodes []ClusterNode
		}
	}

	graphqlErr := n.client.Query(&q, nil)
	if graphqlErr != nil {
		return ClusterNode{}, graphqlErr
	}

	for _, node := range q.Cluster.Nodes {
		if node.Id == nodeID {
			return node, nil
		}
	}

	return ClusterNode{}, fmt.Errorf("node id not found in cluster")
}

func (n *Nodes) Unregister(nodeID int64, force bool) error {
	var m struct {
		ClusterUnregisterNode struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"clusterUnregisterNode(force: $force, nodeID: $id)"`
	}

	variables := map[string]interface{}{
		"id":    graphql.Int(nodeID),
		"force": graphql.Boolean(false),
	}

	graphqlErr := n.client.Mutate(&m, variables)

	return graphqlErr
}
