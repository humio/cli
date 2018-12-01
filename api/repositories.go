package api

import "github.com/shurcooL/graphql"

type Repositories struct {
	client *Client
}

type Repository struct {
	Name            string
	RetentionDays   int64 `graphql:"timeBasedRetention"`
	RetentionSizeGB int64 `graphql:"storageSizeBasedRetention"`
	SpaceUsed       int64 `graphql:"compressedByteSize"`
}

func (c *Client) Repositories() *Repositories { return &Repositories{client: c} }

func (c *Repositories) Get(name string) (Repository, error) {
	var q struct {
		Repository Repository `graphql:"repository(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	graphqlErr := c.client.Query(&q, variables)

	return q.Repository, graphqlErr
}

type RepoListItem struct {
	Name      string
	SpaceUsed int64 `graphql:"compressedByteSize"`
}

func (c *Repositories) List() ([]RepoListItem, error) {
	var q struct {
		Repositories []RepoListItem `graphql:"repositories"`
	}

	variables := map[string]interface{}{}

	graphqlErr := c.client.Query(&q, variables)

	return q.Repositories, graphqlErr
}
