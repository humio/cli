package api

import (
	"fmt"

	"github.com/shurcooL/graphql"
)

type Repositories struct {
	client *Client
}

type Repository struct {
	Name            string
	RetentionDays   float64 `graphql:"timeBasedRetention"`
	RetentionSizeGB float64 `graphql:"storageSizeBasedRetention"`
	SpaceUsed       int64   `graphql:"compressedByteSize"`
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

type Member struct {
	User                 User `graphql:"user"`
	CanAdministrateUsers bool `graphql:"canAdministrateUsers"`
	CanDeleteData        bool `graphql:"canDeleteData"`
}

func (r *Repositories) AddMember(name, username string, adminRights, deleteRights bool) (Member, error) {
	var mutation struct {
		Result struct {
			Member Member
		} `graphql:"addMember(searchDomainName: $name, username: $username, hasMembershipAdminRights: $adminRights, hasDeletionRights: $deleteRights)"`
	}
	variables := map[string]interface{}{
		"name":         graphql.String(name),
		"username":     graphql.String(username),
		"adminRights":  graphql.Boolean(adminRights),
		"deleteRights": graphql.Boolean(deleteRights),
	}

	graphqlErr := r.client.Mutate(&mutation, variables)

	if graphqlErr != nil {
		// The graphql error message is vague if the user already is a member, so add a hint.
		return mutation.Result.Member, fmt.Errorf("%+v. Does the user already have access to the repo?", graphqlErr)
	}

	return mutation.Result.Member, nil
}
