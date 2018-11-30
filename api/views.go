package api

import "github.com/shurcooL/graphql"

type Views struct {
	client *Client
}

type RolePermission struct {
	Role struct {
		Name string
	}
	View struct {
		Name string
	}
	QueryPrefix string
}

type View struct {
	Name  string
	Roles []RolePermission
}

func (c *Client) Views() *Views { return &Views{client: c} }

func (c *Views) Get(name string) (View, error) {
	var q struct {
		View View `graphql:"searchDomain(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	graphqlErr := c.client.Query(&q, variables)

	return q.View, graphqlErr
}

type ViewListItem struct {
	Name string
}

func (c *Views) List() ([]ViewListItem, error) {
	var q struct {
		View []ViewListItem `graphql:"searchDomains"`
	}

	variables := map[string]interface{}{}

	graphqlErr := c.client.Query(&q, variables)

	return q.View, graphqlErr
}
