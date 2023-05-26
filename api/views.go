package api

import (
	"sort"
	"strings"

	"github.com/shurcooL/graphql"
)

type Views struct {
	client *Client
}

type ViewConnection struct {
	RepoName string
	Filter   string
}

type ViewQueryData struct {
	Name        string
	ID          string
	Description string
	ViewInfo    struct {
		Connections []struct {
			Repository struct{ Name string }
			Filter     string
		}
	} `graphql:"... on View"`
}

type View struct {
	Name        string
	ID          string
	Description string
	Connections []ViewConnection
}

func (c *Client) Views() *Views { return &Views{client: c} }

func (c *Views) Get(name string) (*View, error) {
	var query struct {
		Result ViewQueryData `graphql:"searchDomain(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	err := c.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	connections := make([]ViewConnection, len(query.Result.ViewInfo.Connections))
	for i, data := range query.Result.ViewInfo.Connections {
		connections[i] = ViewConnection{
			RepoName: data.Repository.Name,
			Filter:   data.Filter,
		}
	}

	view := View{
		Name:        query.Result.Name,
		ID:          query.Result.ID,
		Description: query.Result.Description,
		Connections: connections,
	}

	return &view, nil
}

type ViewListItem struct {
	Name string
}

func (c *Views) List() ([]ViewListItem, error) {
	var query struct {
		View []ViewListItem `graphql:"searchDomains"`
	}

	err := c.client.Query(&query, nil)

	sort.Slice(query.View, func(i, j int) bool {
		return strings.ToLower(query.View[i].Name) < strings.ToLower(query.View[j].Name)
	})

	return query.View, err
}

type ViewConnectionInput struct {
	RepositoryName graphql.String `json:"repositoryName"`
	Filter         graphql.String `json:"filter"`
}

func (c *Views) Create(name, description string, connections map[string]string) error {
	var mutation struct {
		CreateView struct {
			Name        string
			Description string
		} `graphql:"createView(name: $name, description: $description, connections: $connections)"`
	}

	var viewConnections []ViewConnectionInput
	for k, v := range connections {
		viewConnections = append(
			viewConnections,
			ViewConnectionInput{
				RepositoryName: graphql.String(k),
				Filter:         graphql.String(v),
			})
	}

	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"description": graphql.String(description),
		"connections": viewConnections,
	}

	return c.client.Mutate(&mutation, variables)
}

func (c *Views) Delete(name, reason string) error {
	var mutation struct {
		DeleteSearchDomain struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"deleteSearchDomain(name: $name, deleteMessage: $reason)"`
	}
	variables := map[string]interface{}{
		"name":   graphql.String(name),
		"reason": graphql.String(reason),
	}

	return c.client.Mutate(&mutation, variables)
}

func (c *Views) UpdateConnections(name string, connections map[string]string) error {
	var mutation struct {
		View struct {
			Name string
		} `graphql:"updateView(viewName: $viewName, connections: $connections)"`
	}

	var viewConnections []ViewConnectionInput
	for k, v := range connections {
		viewConnections = append(
			viewConnections,
			ViewConnectionInput{
				RepositoryName: graphql.String(k),
				Filter:         graphql.String(v),
			})
	}

	variables := map[string]interface{}{
		"viewName":    graphql.String(name),
		"connections": viewConnections,
	}

	return c.client.Mutate(&mutation, variables)
}

func (c *Views) UpdateDescription(name string, description string) error {
	var mutation struct {
		UpdateDescriptionMutation struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateDescriptionForSearchDomain(name: $name, newDescription: $description)"`
	}

	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"description": graphql.String(description),
	}

	return c.client.Mutate(&mutation, variables)
}

func (c *Views) AssignRoleToGroup(viewName, groupID, roleID string) error {
	viewData, err := c.Get(viewName)

	if err != nil {
		return err
	}

	viewID := viewData.ID

	var mutation struct {
		AssignRoleToGroup struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"assignRoleToGroup(input:{viewId: $viewId, groupId: $groupId, roleId: $roleId})"`
	}

	variables := map[string]interface{}{
		"viewId":  graphql.String(viewID),
		"groupId": graphql.String(groupID),
		"roleId":  graphql.String(roleID),
	}

	return c.client.Mutate(&mutation, variables)
}
