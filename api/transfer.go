package api

import (
	"errors"
	"github.com/shurcooL/graphql"
	"time"
)

type Transfer struct {
	client *Client
}

func (c *Client) Transfer() *Transfer { return &Transfer{client: c} }

var ErrManagedGroupDoesNotExist = errors.New("managed export group does not exist")

func (t *Transfer) GetManagedExportGroup() (string, error) {
	var q struct {
		ManagedRolesAndGroupsForExport *struct {
			GroupID string
		} `graphql:"managedRolesAndGroupsForExport"`
	}

	err := t.client.Query(&q, nil)

	if err != nil {
		return "", err
	}

	if q.ManagedRolesAndGroupsForExport == nil {
		return "", ErrManagedGroupDoesNotExist
	}

	return q.ManagedRolesAndGroupsForExport.GroupID, nil
}

func (t *Transfer) CreateManagedExportGroup() (string, error) {
	var mutation struct {
		CreateManagedRolesAndGroupsForExport struct {
			GroupID string
		} `graphql:"createManagedRolesAndGroupsForExport"`
	}

	err := t.client.Mutate(&mutation, nil)
	if err != nil {
		return "", err
	}

	return mutation.CreateManagedRolesAndGroupsForExport.GroupID, nil
}

func (t *Transfer) RemoveManagedExportGroup() error {
	var mutation struct {
		RemoveManagedRolesAndGroupsForExport bool `graphql:"removeManagedRolesAndGroupsForExport"`
	}

	err := t.client.Mutate(&mutation, nil)
	return err
}

type TransferJob struct {
	ID                       string
	SourceClusterURL         string
	Dataspaces               []string
	MaximumParallelDownloads int
	CompletedAt              *time.Time
	CancelledAt              *time.Time
}

func (t *Transfer) ListTransferJobs() ([]TransferJob, error) {
	var q struct {
		TransferJobs []TransferJob `graphql:"transferJobs"`
	}

	err := t.client.Query(&q, nil)
	if err != nil {
		return nil, err
	}

	return q.TransferJobs, nil
}

type AddTransferJobResponse struct {
	ID string
}

func (t *Transfer) AddTransferJob(sourceClusterURL string, sourceClusterToken string, destinationOrganizationID string, dataspaces []string, maximumParallelDownloads int, setTargetAsNewMaster bool) (AddTransferJobResponse, error) {
	var mutation struct {
		AddTransferJob AddTransferJobResponse `graphql:"addTransferJob(input: {sourceClusterUrl: $sourceClusterUrl, sourceClusterToken: $sourceClusterToken, destinationOrganizationId: $destinationOrganizationId, dataspaces: $dataspaces, maximumParallelDownloads: $maximumParallelDownloads, setTargetClusterAsNewMaster: $setTargetClusterAsNewMaster})"`
	}

	ds := make([]graphql.String, len(dataspaces))
	for i := range dataspaces {
		ds[i] = graphql.String(dataspaces[i])
	}

	variables := map[string]interface{}{
		"sourceClusterUrl":            graphql.String(sourceClusterURL),
		"sourceClusterToken":          graphql.String(sourceClusterToken),
		"destinationOrganizationId":   graphql.String(destinationOrganizationID),
		"dataspaces":                  ds,
		"setTargetClusterAsNewMaster": graphql.Boolean(setTargetAsNewMaster),
	}

	if maximumParallelDownloads > 0 {
		variables["maximumParallelDownloads"] = graphql.Int(maximumParallelDownloads)
	} else {
		variables["maximumParallelDownloads"] = (*graphql.Int)(nil)
	}

	err := t.client.Mutate(&mutation, variables)

	return mutation.AddTransferJob, err
}

func (t *Transfer) CancelTransferJob(transferJobID string) (TransferJob, error) {
	var mutation struct {
		CancelledTransferJob TransferJob `graphql:"cancelTransferJob(transferJobId: $transferJobId)"`
	}

	variables := map[string]interface{}{
		"transferJobId": graphql.String(transferJobID),
	}

	err := t.client.Mutate(&mutation, variables)

	return mutation.CancelledTransferJob, err
}

type TransferJobStatus struct {
	TotalSegments       int
	TransferredSegments int
	Running             bool
	Error               string
	Status              string
	StatusLine          string
}

func (t *Transfer) GetTransferJobStatus(transferJobID string) (TransferJobStatus, error) {
	var q struct {
		TransferJobStatus TransferJobStatus `graphql:"transferJobStatus(transferJobId: $transferJobId)"`
	}

	variables := map[string]interface{}{
		"transferJobId": graphql.String(transferJobID),
	}

	err := t.client.Query(&q, variables)

	return q.TransferJobStatus, err
}
