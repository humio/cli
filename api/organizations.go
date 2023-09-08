package api

import graphql "github.com/cli/shurcooL-graphql"

type Organizations struct {
	client *Client
}

type Organization struct {
	ID          string
	Name        string
	Description *string
}

func (c *Client) Organizations() *Organizations { return &Organizations{client: c} }

func (o *Organizations) CreateOrganization(name string) (Organization, error) {
	var mutation struct {
		CreateOrganization Organization `graphql:"createEmptyOrganization(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	err := o.client.Mutate(&mutation, variables)

	return mutation.CreateOrganization, err
}
