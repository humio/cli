package api

import "github.com/shurcooL/graphql"

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
	var m struct {
		CreateOrganization Organization `graphql:"createEmptyOrganization(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	err := o.client.Mutate(&m, variables)

	return m.CreateOrganization, err
}
