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

	return q.Repository, graphqlErr
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
