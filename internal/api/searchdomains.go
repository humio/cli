package api

import (
	"context"
	"sort"
	"strings"

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

func (s *Client) SearchDomains() *SearchDomains { return &SearchDomains{client: s} }

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

type SearchDomainListItem struct {
	Name            string
	Typename        *string
	AutomaticSearch bool
}

func (s *SearchDomains) List() ([]SearchDomainListItem, error) {
	resp, err := humiographql.ListSearchDomains(context.Background(), s.client)
	if err != nil {
		return nil, err
	}

	searchDomains := resp.GetSearchDomains()
	searchDomainList := []SearchDomainListItem{}
	for _, searchDomain := range searchDomains {
		searchDomainList = append(searchDomainList, SearchDomainListItem{
			Name:            searchDomain.GetName(),
			Typename:        searchDomain.GetTypename(),
			AutomaticSearch: searchDomain.GetAutomaticSearch(),
		})
	}

	sort.Slice(searchDomainList, func(i, j int) bool {
		return strings.ToLower(searchDomainList[i].Name) < strings.ToLower(searchDomainList[j].Name)
	})
	return searchDomainList, nil
}
