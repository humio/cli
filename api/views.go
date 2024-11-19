package api

import (
	"sort"
	"strings"

	"github.com/humio/cli/api/internal/humiographql"

	graphql "github.com/cli/shurcooL-graphql"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Views struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type ViewConnection struct {
	RepoName string
	Filter   string
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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
	Typename graphql.String `graphql:"__typename"`
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type View struct {
	Name            string
	Description     string
	Connections     []ViewConnection
	AutomaticSearch bool
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Views() *Views { return &Views{client: c} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Views) Get(name string) (*View, error) {
	var query struct {
		Result ViewQueryData `graphql:"searchDomain(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	err := c.client.Query(&query, variables)
	if err != nil {
		return nil, ViewNotFound(name)
	}
	if query.Result.Typename != "View" {
		return nil, ViewNotFound("name")
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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type ViewListItem struct {
	Name            string
	Typename        string `graphql:"__typename"`
	AutomaticSearch bool
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Views) List() ([]ViewListItem, error) {
	var query struct {
		View []ViewListItem `graphql:"searchDomains"`
	}

	err := c.client.Query(&query, nil)

	viewsList := []ViewListItem{}
	for k, v := range query.View {
		if v.Typename == string(humiographql.SearchDomainTypeView) {
			viewsList = append(viewsList, query.View[k])
		}
	}

	sort.Slice(viewsList, func(i, j int) bool {
		return strings.ToLower(viewsList[i].Name) < strings.ToLower(viewsList[j].Name)
	})

	return viewsList, err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type ViewConnectionInput struct {
	RepositoryName graphql.String `json:"repositoryName"`
	Filter         graphql.String `json:"filter"`
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Views) Delete(name, reason string) error {
	_, err := c.Get(name)
	if err != nil {
		return err
	}

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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Views) UpdateDescription(name string, description string) error {
	_, err := c.Get(name)
	if err != nil {
		return err
	}

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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Views) UpdateAutomaticSearch(name string, automaticSearch bool) error {
	_, err := c.Get(name)
	if err != nil {
		return err
	}

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
