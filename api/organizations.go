package api

import graphql "github.com/cli/shurcooL-graphql"

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Organizations struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Organization struct {
	ID          string
	Name        string
	Description *string
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Organizations() *Organizations { return &Organizations{client: c} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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
