package api

import (
	"context"
	"strconv"

	"github.com/humio/cli/internal/api/humiographql"
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
	Zone *string
}

type IngestPartition struct {
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
}

func (c *Client) Clusters() *Clusters { return &Clusters{client: c} }

func (c *Clusters) Get() (Cluster, error) {
	resp, err := humiographql.GetCluster(context.Background(), c.client)
	if err != nil {
		return Cluster{}, err
	}

	respCluster := resp.GetCluster()
	respClusterNodes := respCluster.GetNodes()
	clusterNodes := make([]ClusterNode, len(respClusterNodes))
	for idx, node := range respClusterNodes {
		clusterNodes[idx] = ClusterNode{
			Id:                      node.GetId(),
			Name:                    node.GetName(),
			Uri:                     node.GetUri(),
			Uuid:                    node.GetUuid(),
			ClusterInfoAgeSeconds:   node.GetClusterInfoAgeSeconds(),
			InboundSegmentSize:      node.GetInboundSegmentSize(),
			OutboundSegmentSize:     node.GetOutboundSegmentSize(),
			CanBeSafelyUnregistered: node.GetCanBeSafelyUnregistered(),
			CurrentSize:             node.GetCurrentSize(),
			PrimarySize:             node.GetPrimarySize(),
			SecondarySize:           node.GetSecondarySize(),
			TotalSizeOfPrimary:      node.GetTotalSizeOfPrimary(),
			TotalSizeOfSecondary:    node.GetTotalSizeOfSecondary(),
			FreeOnPrimary:           node.GetFreeOnPrimary(),
			FreeOnSecondary:         node.GetFreeOnSecondary(),
			WipSize:                 node.GetWipSize(),
			TargetSize:              node.GetTargetSize(),
			SolitarySegmentSize:     node.GetSolitarySegmentSize(),
			IsAvailable:             node.GetIsAvailable(),
			LastHeartbeat:           node.GetLastHeartbeat().String(),
			Zone:                    node.GetZone(),
		}
	}
	respIngestPartitions := respCluster.GetIngestPartitions()
	ingestPartitions := make([]IngestPartition, len(respIngestPartitions))
	for idx, partition := range respIngestPartitions {
		ingestPartitions[idx] = IngestPartition{
			Id:      partition.GetId(),
			NodeIds: partition.GetNodeIds(),
		}
	}

	return Cluster{
		Nodes:                               clusterNodes,
		ClusterInfoAgeSeconds:               respCluster.GetClusterInfoAgeSeconds(),
		UnderReplicatedSegmentSize:          respCluster.GetUnderReplicatedSegmentSize(),
		OverReplicatedSegmentSize:           respCluster.GetOverReplicatedSegmentSize(),
		MissingSegmentSize:                  respCluster.GetMissingSegmentSize(),
		ProperlyReplicatedSegmentSize:       respCluster.GetProperlyReplicatedSegmentSize(),
		TargetUnderReplicatedSegmentSize:    respCluster.GetTargetUnderReplicatedSegmentSize(),
		TargetOverReplicatedSegmentSize:     respCluster.GetTargetOverReplicatedSegmentSize(),
		TargetMissingSegmentSize:            respCluster.GetTargetMissingSegmentSize(),
		TargetProperlyReplicatedSegmentSize: respCluster.GetTargetProperlyReplicatedSegmentSize(),
		IngestPartitions:                    ingestPartitions,
	}, nil
}

type ClusterNodes struct {
	client *Client
}

func (c *Client) ClusterNodes() *ClusterNodes { return &ClusterNodes{client: c} }

func (n *ClusterNodes) List() ([]ClusterNode, error) {
	resp, err := humiographql.GetCluster(context.Background(), n.client)
	if err != nil {
		return nil, err
	}

	respCluster := resp.GetCluster()
	respClusterNodes := respCluster.GetNodes()
	clusterNodes := make([]ClusterNode, len(respClusterNodes))
	for idx, node := range respClusterNodes {
		clusterNodes[idx] = ClusterNode{
			Id:                      node.GetId(),
			Name:                    node.GetName(),
			Uri:                     node.GetUri(),
			Uuid:                    node.GetUuid(),
			ClusterInfoAgeSeconds:   node.GetClusterInfoAgeSeconds(),
			InboundSegmentSize:      node.GetInboundSegmentSize(),
			OutboundSegmentSize:     node.GetOutboundSegmentSize(),
			CanBeSafelyUnregistered: node.GetCanBeSafelyUnregistered(),
			CurrentSize:             node.GetCurrentSize(),
			PrimarySize:             node.GetPrimarySize(),
			SecondarySize:           node.GetSecondarySize(),
			TotalSizeOfPrimary:      node.GetTotalSizeOfPrimary(),
			TotalSizeOfSecondary:    node.GetTotalSizeOfSecondary(),
			FreeOnPrimary:           node.GetFreeOnPrimary(),
			FreeOnSecondary:         node.GetFreeOnSecondary(),
			WipSize:                 node.GetWipSize(),
			TargetSize:              node.GetTargetSize(),
			SolitarySegmentSize:     node.GetSolitarySegmentSize(),
			IsAvailable:             node.GetIsAvailable(),
			LastHeartbeat:           node.GetLastHeartbeat().String(),
			Zone:                    node.GetZone(),
		}
	}

	return clusterNodes, nil
}

func (n *ClusterNodes) Get(nodeID int) (ClusterNode, error) {
	resp, err := humiographql.GetCluster(context.Background(), n.client)
	if err != nil {
		return ClusterNode{}, err
	}

	respCluster := resp.GetCluster()
	respClusterNodes := respCluster.GetNodes()
	for _, node := range respClusterNodes {
		if node.Id == nodeID {
			return ClusterNode{
				Id:                      node.GetId(),
				Name:                    node.GetName(),
				Uri:                     node.GetUri(),
				Uuid:                    node.GetUuid(),
				ClusterInfoAgeSeconds:   node.GetClusterInfoAgeSeconds(),
				InboundSegmentSize:      node.GetInboundSegmentSize(),
				OutboundSegmentSize:     node.GetOutboundSegmentSize(),
				CanBeSafelyUnregistered: node.GetCanBeSafelyUnregistered(),
				CurrentSize:             node.GetCurrentSize(),
				PrimarySize:             node.GetPrimarySize(),
				SecondarySize:           node.GetSecondarySize(),
				TotalSizeOfPrimary:      node.GetTotalSizeOfPrimary(),
				TotalSizeOfSecondary:    node.GetTotalSizeOfSecondary(),
				FreeOnPrimary:           node.GetFreeOnPrimary(),
				FreeOnSecondary:         node.GetFreeOnSecondary(),
				WipSize:                 node.GetWipSize(),
				TargetSize:              node.GetTargetSize(),
				SolitarySegmentSize:     node.GetSolitarySegmentSize(),
				IsAvailable:             node.GetIsAvailable(),
				LastHeartbeat:           node.GetLastHeartbeat().String(),
				Zone:                    node.GetZone(),
			}, nil
		}
	}

	return ClusterNode{}, ClusterNodeNotFound(strconv.Itoa(nodeID))
}

func (n *ClusterNodes) Unregister(nodeID int, force bool) error {
	_, err := humiographql.UnregisterClusterNode(context.Background(), n.client, nodeID, force)
	if err != nil {
		return err
	}
	return nil
}
