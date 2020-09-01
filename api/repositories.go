package api

import (
	"fmt"

	"github.com/shurcooL/graphql"
)

type Repositories struct {
	client *Client
}

type Repository struct {
	Name                   string
	Description            string
	RetentionDays          float64 `graphql:"timeBasedRetention"`
	IngestRetentionSizeGB  float64 `graphql:"ingestSizeBasedRetention"`
	StorageRetentionSizeGB float64 `graphql:"storageSizeBasedRetention"`
	SpaceUsed              int64   `graphql:"compressedByteSize"`
}

func (c *Client) Repositories() *Repositories { return &Repositories{client: c} }

func (r *Repositories) Get(name string) (Repository, error) {
	var q struct {
		Repository Repository `graphql:"repository(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	graphqlErr := r.client.Query(&q, variables)

	if graphqlErr != nil {
		// The graphql error message is vague if the repo already exists, so add a hint.
		return q.Repository, fmt.Errorf("%+v. Does the repo already exist?", graphqlErr)
	}

	return q.Repository, nil
}

type RepoListItem struct {
	Name      string
	SpaceUsed int64 `graphql:"compressedByteSize"`
}

func (r *Repositories) List() ([]RepoListItem, error) {
	var q struct {
		Repositories []RepoListItem `graphql:"repositories"`
	}

	graphqlErr := r.client.Query(&q, nil)

	return q.Repositories, graphqlErr
}

func (r *Repositories) Create(name string) error {
	var m struct {
		CreateRepository struct {
			Repository Repository
		} `graphql:"createRepository(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	graphqlErr := r.client.Mutate(&m, variables)

	if graphqlErr != nil {
		// The graphql error message is vague if the repo already exists, so add a hint.
		return fmt.Errorf("%+v. Does the repo already exist?", graphqlErr)
	}

	return nil
}

func (r *Repositories) Delete(name, reason string, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}
	safeToDelete := allowDataDeletion || existingRepo.SpaceUsed == 0
	if !safeToDelete {
		return fmt.Errorf("repository contains data and data deletion not allowed")
	}

	var m struct {
		CreateRepository struct {
			Type string `graphql:"__typename"`
		} `graphql:"deleteSearchDomain(name: $name, deleteMessage: $reason)"`
	}
	variables := map[string]interface{}{
		"name":   graphql.String(name),
		"reason": graphql.String(reason),
	}

	err = r.client.Mutate(&m, variables)

	if err != nil {
		return err
	}

	return nil
}

type DefaultGroupEnum string

func (r *Repositories) AddUser(name, username, group string) error {
	var mutation struct {
		UpdateDefaultGroupMembershipsMutation struct {
			ClientMutationId string
		} `graphql:"updateDefaultGroupMemberships(input: {viewName: $name, userName: $username, groups: [$group]})"`
	}
	variables := map[string]interface{}{
		"name":     graphql.String(name),
		"username": graphql.String(username),
		"group":    DefaultGroupEnum(group),
	}

	graphqlErr := r.client.Mutate(&mutation, variables)

	if graphqlErr != nil {
		// The graphql error message is vague if the user already is a member, so add a hint.
		return fmt.Errorf("%+v. Does the user already have access to the repo?", graphqlErr)
	}

	return nil
}

func (r *Repositories) UpdateTimeBasedRetention(name string, retentionInDays float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}
	safeToDelete := allowDataDeletion || existingRepo.SpaceUsed == 0

	var m struct {
		UpdateRetention struct {
			Type string `graphql:"__typename"`
		} `graphql:"updateRetention(repositoryName: $name, timeBasedRetention: $retentionInDays)"`
	}
	variables := map[string]interface{}{
		"name":            graphql.String(name),
		"retentionInDays": (*graphql.Float)(nil),
	}
	if retentionInDays > 0 {
		if retentionInDays < existingRepo.RetentionDays || existingRepo.RetentionDays == 0 {
			if !safeToDelete {
				return fmt.Errorf("repository contains data and data deletion not allowed")
			}
		}
		variables["retentionInDays"] = graphql.Float(retentionInDays)
	}

	err = r.client.Mutate(&m, variables)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repositories) UpdateStorageBasedRetention(name string, storageInGB float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}
	safeToDelete := allowDataDeletion || existingRepo.SpaceUsed == 0

	var m struct {
		UpdateRetention struct {
			Type string `graphql:"__typename"`
		} `graphql:"updateRetention(repositoryName: $name, storageSizeBasedRetention: $storageInGB)"`
	}
	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"storageInGB": (*graphql.Float)(nil),
	}
	if storageInGB > 0 {
		if storageInGB < existingRepo.StorageRetentionSizeGB || existingRepo.StorageRetentionSizeGB == 0 {
			if !safeToDelete {
				return fmt.Errorf("repository contains data and data deletion not allowed")
			}
		}
		variables["storageInGB"] = graphql.Float(storageInGB)
	}

	err = r.client.Mutate(&m, variables)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repositories) UpdateIngestBasedRetention(name string, ingestInGB float64, allowDataDeletion bool) error {
	existingRepo, err := r.Get(name)
	if err != nil {
		return err
	}
	safeToDelete := allowDataDeletion || existingRepo.SpaceUsed == 0

	var m struct {
		UpdateRetention struct {
			Type string `graphql:"__typename"`
		} `graphql:"updateRetention(repositoryName: $name, ingestSizeBasedRetention: $ingestInGB)"`
	}
	variables := map[string]interface{}{
		"name":       graphql.String(name),
		"ingestInGB": (*graphql.Float)(nil),
	}
	if ingestInGB > 0 {
		if ingestInGB < existingRepo.IngestRetentionSizeGB || existingRepo.IngestRetentionSizeGB == 0 {
			if !safeToDelete {
				return fmt.Errorf("repository contains data and data deletion not allowed")
			}
		}
		variables["ingestInGB"] = graphql.Float(ingestInGB)
	}

	err = r.client.Mutate(&m, variables)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repositories) UpdateDescription(name, description string) error {
	var m struct {
		UpdateDescription struct {
			Type string `graphql:"__typename"`
		} `graphql:"updateDescriptionForSearchDomain(name: $name, newDescription: $description)"`
	}

	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"description": graphql.String(description),
	}

	err := r.client.Mutate(&m, variables)

	if err != nil {
		return err
	}

	return nil
}
