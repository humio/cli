package api

import (
	"fmt"
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
	Filter graphql.String `json:"filter"`
}

func (c *Views) Create(name, description string, connections map[string]string) error {
	var m struct {
		CreateView struct {
			Name string
			Description string
		} `graphql:"createView(name: $name, description: $description, connections: $connections)"`
	}
                                                                                                
	viewConnections := make([]ViewConnectionInput, len(connections))
	i := 0
	for k, v := range connections {
		viewConnections[i] = ViewConnectionInput{
			RepositoryName: graphql.String(k),
			Filter: graphql.String(v),
		}

		i++
	}

	variables := map[string]interface{} {
		"name": graphql.String(name),
		"description": graphql.String(description),
		"connections": viewConnections,
	}

	graphqlErr := c.client.Mutate(&m, variables)

	if graphqlErr != nil {
		// The graphql error message is vague if the repo already exists, so add a hint.
		return fmt.Errorf("%+v. Does the view already exist?", graphqlErr)
	}

	return nil
}

func (c *Views) Delete(name, reason string) error {
	var m struct {
		DeleteSearchDomain struct {
			ClientMutationId string
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