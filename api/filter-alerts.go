package api

import (
	"fmt"

	graphql "github.com/cli/shurcooL-graphql"
)

type FilterAlert struct {
	ID                 string             `graphql:"id"             yaml:"-"                       json:"id"`
	Name               string             `graphql:"name"           yaml:"name"                    json:"name"`
	Description        string             `graphql:"description"    yaml:"description,omitempty"   json:"description,omitempty"`
	QueryString        string             `graphql:"queryString"    yaml:"queryString"             json:"queryString"`
	ActionNames        []string           `graphql:"actionNames"    yaml:"actionNames"             json:"actionNames"`
	Labels             []string           `graphql:"labels"         yaml:"labels"                  json:"labels"`
	Enabled            bool               `graphql:"enabled"        yaml:"enabled"                 json:"enabled"`
	QueryOwnershipType QueryOwnershipType `graphql:"queryOwnership" yaml:"queryOwnershipType"      json:"queryOwnershipType"`
	RunAsUserID        string             `graphql:"runAsUserId"    yaml:"runAsUserId,omitempty"   json:"runAsUserId,omitempty"`
}

type FilterAlerts struct {
	client *Client
}

func (c *Client) FilterAlerts() *FilterAlerts { return &FilterAlerts{client: c} }

func (fa *FilterAlerts) List(viewName string) ([]FilterAlert, error) {
	var query struct {
		SearchDomain struct {
			FilterAlerts []FilterAlert `graphql:"filterAlerts"`
		} `graphql:"searchDomain(name: $viewName)"`
	}

	variables := map[string]any{
		"viewName": graphql.String(viewName),
	}
	err := fa.client.Query(&query, variables)
	return query.SearchDomain.FilterAlerts, err
}

func (fa *FilterAlerts) Update(viewName string, updatedFilterAlert UpdateFilterAlert) (*FilterAlert, error) {
	var mutation struct {
		FilterAlert `graphql:"updateFilterAlert(input: $input)"`
	}

	variables := map[string]any{
		"input": updatedFilterAlert,
	}
	err := fa.client.Mutate(&mutation, variables)
	return &mutation.FilterAlert, err
}

func (fa *FilterAlerts) Create(viewName string, newFilterAlert CreateFilterAlert) (*FilterAlert, error) {
	var mutation struct {
		FilterAlert FilterAlert `graphql:"createFilterAlert(input: $input)"`
	}

	variables := map[string]any{
		"input": newFilterAlert,
	}
	err := fa.client.Mutate(&mutation, variables)
	return &mutation.FilterAlert, err
}

func (fa *FilterAlerts) Delete(viewName, filterAlertID string) error {
	if filterAlertID == "" {
		return fmt.Errorf("filterAlertID is empty")
	}

	var mutation struct {
		DidDelete bool `graphql:"deleteFilterAlert(input: { viewName: $viewName, id: $id })"`
	}

	variables := map[string]any{
		"viewName": RepoOrViewName(viewName),
		"id":       graphql.String(filterAlertID),
	}

	err := fa.client.Mutate(&mutation, variables)

	if !mutation.DidDelete {
		return fmt.Errorf("unable to remove filter alert in repo/view '%s' with id '%s'", viewName, filterAlertID)
	}

	return err
}

func (fa *FilterAlerts) Get(viewName string, filterAlertID string) (*FilterAlert, error) {
	var query struct {
		SearchDomain struct {
			FilterAlert FilterAlert `graphql:"filterAlert(id: $filterAlertId)"`
		} `graphql:"searchDomain(name: $viewName) "`
	}

	variables := map[string]any{
		"viewName":      graphql.String(viewName),
		"filterAlertId": graphql.String(filterAlertID),
	}

	err := fa.client.Query(&query, variables)
	return &query.SearchDomain.FilterAlert, err
}
