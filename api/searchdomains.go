package api

import (
	"sort"
	"strings"

	graphql "github.com/cli/shurcooL-graphql"
)

type SearchDomains struct {
	client *Client
}

type SearchDomainsQueryData struct {
	Name            string
	Description     string
	AutomaticSearch bool
	Typename        graphql.String `graphql:"__typename"`
}

type SearchDomain struct {
	Name            string
	Description     string
	AutomaticSearch bool
}

func (s *Client) SearchDomains() *SearchDomains { return &SearchDomains{client: s} }

func (s *SearchDomains) Get(name string) (*SearchDomain, error) {
	var query struct {
		Result SearchDomainsQueryData `graphql:"searchDomain(name: $name)"`
	}

	variables := map[string]interface{}{
		"name": graphql.String(name),
	}

	err := s.client.Query(&query, variables)
	if err != nil {
		return nil, SearchDomainNotFound(name)
	}

	searchDomain := SearchDomain{
		Name:            query.Result.Name,
		Description:     query.Result.Description,
		AutomaticSearch: query.Result.AutomaticSearch,
	}

	return &searchDomain, nil
}

type SearchDomainListItem struct {
	Name            string
	Typename        string `graphql:"__typename"`
	AutomaticSearch bool
}

func (s *SearchDomains) List() ([]SearchDomainListItem, error) {
	var query struct {
		SearchDomain []SearchDomainListItem `graphql:"searchDomains"`
	}

	err := s.client.Query(&query, nil)

	sort.Slice(query.SearchDomain, func(i, j int) bool {
		return strings.ToLower(query.SearchDomain[i].Name) < strings.ToLower(query.SearchDomain[j].Name)
	})

	return query.SearchDomain, err
}
