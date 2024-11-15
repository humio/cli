package api

import (
	"sort"
	"strings"

	graphql "github.com/cli/shurcooL-graphql"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type SearchDomains struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type SearchDomainsQueryData struct {
	Name            string
	Description     string
	AutomaticSearch bool
	Typename        graphql.String `graphql:"__typename"`
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type SearchDomain struct {
	Name            string
	Description     string
	AutomaticSearch bool
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (s *Client) SearchDomains() *SearchDomains { return &SearchDomains{client: s} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type SearchDomainListItem struct {
	Name            string
	Typename        string `graphql:"__typename"`
	AutomaticSearch bool
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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
