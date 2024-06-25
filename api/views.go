package api

import (
	"sort"
	"strings"

	graphql "github.com/cli/shurcooL-graphql"
)

type Views struct {
	client *Client
}

type ViewConnection struct {
	RepoName string
	Filter   string
}

type ViewQueryData struct {
	Name            string
	Description     string
	AutomaticSearch bool
	ViewInfo        struct {
		Connections []struct {
			Repository struct{ Name string }
			Filter     string
		}
	} `graphql:"... on View"`
}

type View struct {
	Name            string
	Description     string
	Connections     []ViewConnection
	AutomaticSearch bool
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
		Name:            query.Result.Name,
		Description:     query.Result.Description,
		Connections:     connections,
		AutomaticSearch: query.Result.AutomaticSearch,
	}

	return &view, nil
}

type ViewListItem struct {
	Name            string
	Typename        string `graphql:"__typename"`
	AutomaticSearch bool
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

func (c *Views) Create(name, description string, connections []ViewConnectionInput) error {
	var mutation struct {
		CreateView struct {
			Name        string
			Description string
		} `graphql:"createView(name: $name, description: $description, connections: $connections)"`
	}

	variables := map[string]interface{}{
		"name":        graphql.String(name),
		"description": graphql.String(description),
		"connections": connections,
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

func (c *Views) UpdateConnections(name string, connections []ViewConnectionInput) error {
	var mutation struct {
		View struct {
			Name string
		} `graphql:"updateView(viewName: $viewName, connections: $connections)"`
	}

	variables := map[string]interface{}{
		"viewName":    graphql.String(name),
		"connections": connections,
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

func (c *Views) UpdateAutomaticSearch(name string, automaticSearch bool) error {
	var mutation struct {
		SetAutomaticSearching struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"setAutomaticSearching(name: $name, automaticSearch: $automaticSearch)"`
	}

	variables := map[string]interface{}{
		"name":            graphql.String(name),
		"automaticSearch": graphql.Boolean(automaticSearch),
	}

	return c.client.Mutate(&mutation, variables)
}
