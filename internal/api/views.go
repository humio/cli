package api

import (
	"context"
	"sort"
	"strings"

	"github.com/humio/cli/internal/api/humiographql"
)

type Views struct {
	client *Client
}

type ViewConnection struct {
	RepoName string
	Filter   string
}

type View struct {
	Name            string
	Description     string
	Connections     []ViewConnection
	AutomaticSearch bool
}

func (c *Client) Views() *Views { return &Views{client: c} }

func (c *Views) Get(name string) (*View, error) {
	resp, err := humiographql.GetSearchDomain(context.Background(), c.client, name)
	if err != nil {
		return nil, ViewNotFound(name)
	}

	searchDomain := resp.GetSearchDomain()

	switch v := searchDomain.(type) {
	case *humiographql.GetSearchDomainSearchDomainView:
		connections := make([]ViewConnection, len(v.GetConnections()))
		for i, data := range v.GetConnections() {
			connections[i] = ViewConnection{
				RepoName: data.Repository.Name,
				Filter:   data.Filter,
			}
		}
		description := ""
		if searchDomain.GetDescription() != nil {
			description = *searchDomain.GetDescription()
		}
		return &View{
			Name:            searchDomain.GetName(),
			Description:     description,
			Connections:     connections,
			AutomaticSearch: searchDomain.GetAutomaticSearch(),
		}, nil
	default:
		return nil, ViewNotFound(name)
	}
}

type ViewListItem struct {
	Name            string
	Typename        string
	AutomaticSearch bool
}

func (c *Views) List() ([]ViewListItem, error) {
	resp, err := humiographql.ListSearchDomains(context.Background(), c.client)
	if err != nil {
		return nil, err
	}

	searchDomains := resp.GetSearchDomains()
	viewsList := []ViewListItem{}
	for _, searchDomain := range searchDomains {
		switch v := searchDomain.(type) {
		case *humiographql.ListSearchDomainsSearchDomainsView:
			typename := ""
			if v.GetTypename() != nil {
				typename = *v.GetTypename()
			}
			viewsList = append(viewsList, ViewListItem{
				Name:            v.GetName(),
				Typename:        typename,
				AutomaticSearch: v.GetAutomaticSearch(),
			})
		default:
			// ignore
		}
	}

	sort.Slice(viewsList, func(i, j int) bool {
		return strings.ToLower(viewsList[i].Name) < strings.ToLower(viewsList[j].Name)
	})
	return viewsList, nil
}

type ViewConnectionInput struct {
	RepositoryName string
	Filter         string
}

func (c *Views) Create(name, description string, connections []ViewConnectionInput) error {
	createDescription := ""
	if description != "" {
		createDescription = description
	}
	internalConnType := make([]humiographql.ViewConnectionInput, len(connections))
	for i := range connections {
		internalConnType[i] = humiographql.ViewConnectionInput{
			RepositoryName: connections[i].RepositoryName,
			Filter:         connections[i].Filter,
		}
	}
	_, err := humiographql.CreateView(context.Background(), c.client, name, &createDescription, internalConnType)
	return err
}

func (c *Views) Delete(name, reason string) error {
	_, err := c.Get(name)
	if err != nil {
		return err
	}

	_, err = humiographql.DeleteSearchDomain(context.Background(), c.client, name, reason)
	return err
}

func (c *Views) UpdateConnections(name string, connections []ViewConnectionInput) error {
	internalConnType := make([]humiographql.ViewConnectionInput, len(connections))
	for i := range connections {
		internalConnType[i] = humiographql.ViewConnectionInput{
			RepositoryName: connections[i].RepositoryName,
			Filter:         connections[i].Filter,
		}
	}
	_, err := humiographql.UpdateViewConnections(context.Background(), c.client, name, internalConnType)
	return err
}

func (c *Views) UpdateDescription(name string, description string) error {
	_, err := c.Get(name)
	if err != nil {
		return err
	}

	_, err = humiographql.UpdateDescriptionForSearchDomain(context.Background(), c.client, name, description)
	return err
}

func (c *Views) UpdateAutomaticSearch(name string, automaticSearch bool) error {
	_, err := c.Get(name)
	if err != nil {
		return err
	}

	_, err = humiographql.SetAutomaticSearching(context.Background(), c.client, name, automaticSearch)
	return err
}
