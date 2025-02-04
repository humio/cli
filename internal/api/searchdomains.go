package api

import (
	"context"

	"github.com/humio/cli/internal/api/humiographql"
)

type SearchDomains struct {
	client *Client
}

type SearchDomain struct {
	Name            string
	Description     *string
	AutomaticSearch bool
}

func (c *Client) SearchDomains() *SearchDomains { return &SearchDomains{client: c} }

func (s *SearchDomains) Get(name string) (*SearchDomain, error) {
	resp, err := humiographql.GetSearchDomain(context.Background(), s.client, name)
	if err != nil {
		return nil, SearchDomainNotFound(name)
	}

	searchDomain := resp.GetSearchDomain()
	return &SearchDomain{
		Name:            searchDomain.GetName(),
		Description:     searchDomain.GetDescription(),
		AutomaticSearch: searchDomain.GetAutomaticSearch(),
	}, nil
}
