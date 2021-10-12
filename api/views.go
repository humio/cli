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
	Name     string
	ViewInfo struct {
		Connections []struct {
			Repository struct{ Name string }
			Filter     string
		}
	} `graphql:"... on View"`
}

type View struct {
	Name        string
	Connections []ViewConnection
}

func (c *Client) Views() *Views { return &Views{client: c} }

func (c *Views) Get(name string) (*View, error) {
	var q struct {
		Result ViewQueryData `graphql:"searchDomain(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	graphqlErr := c.client.Query(&q, variables)

	if graphqlErr != nil {
		return nil, graphqlErr
	}

	connections := make([]ViewConnection, len(q.Result.ViewInfo.Connections))

	for i, data := range q.Result.ViewInfo.Connections {
		connections[i] = ViewConnection{
			RepoName: data.Repository.Name,
			Filter:   data.Filter,
		}
	}

	view := View{
		Name:        q.Result.Name,
		Connections: connections,
	}

	return &view, nil
}

type ViewListItem struct {
	Name string
}

func (c *Views) List() ([]ViewListItem, error) {
	var q struct {
		View []ViewListItem `graphql:"searchDomains"`
	}

	graphqlErr := c.client.Query(&q, nil)

	sort.Slice(q.View, func(i, j int) bool {
		return strings.ToLower(q.View[i].Name) < strings.ToLower(q.View[j].Name)
	})

	return q.View, graphqlErr
}

type ViewConnectionInput struct {
	RepositoryName graphql.String `json:"repositoryName"`
	Filter         graphql.String `json:"filter"`
}

func (c *Views) Create(name, description string, connections map[string]string) error {
	var m struct {
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

	err := c.client.Mutate(&m, variables)

	if err != nil {
		return err
	}

	return nil
}

func (c *Views) Delete(name, reason string) error {
	var m struct {
		DeleteSearchDomain struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"deleteSearchDomain(name: $name, deleteMessage: $reason)"`
	}
	variables := map[string]interface{}{
		"name":   graphql.String(name),
		"reason": graphql.String(reason),
	}

	err := c.client.Mutate(&m, variables)

	if err != nil {
		return err
	}

	return nil
}

func (c *Views) UpdateConnections(name string, connections map[string]string) error {
	var m struct {
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

	err := c.client.Mutate(&m, variables)

	if err != nil {
		return err
	}

	return nil
}

func (c *Views) UpdateDescription(name string, description string) error {
	var m struct {
		UpdateDescriptionMutation struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateDescriptionForSearchDomain(name: $name, newDescription: $description)"`
	}

	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"description": graphql.String(description),
	}

	err := c.client.Mutate(&m, variables)

	if err != nil {
		return err
	}

	return nil
}
