package api

import (
	"fmt"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/humio/cli/api/internal/humiographql"
)

type FilterAlert struct {
	ID                  string   `graphql:"id"                  yaml:"-"                             json:"id"`
	Name                string   `graphql:"name"                yaml:"name"                          json:"name"`
	Description         string   `graphql:"description"         yaml:"description,omitempty"         json:"description,omitempty"`
	QueryString         string   `graphql:"queryString"         yaml:"queryString"                   json:"queryString"`
	ActionNames         []string `graphql:"actionNames"         yaml:"actionNames"                   json:"actionNames"`
	Labels              []string `graphql:"labels"              yaml:"labels"                        json:"labels"`
	Enabled             bool     `graphql:"enabled"             yaml:"enabled"                       json:"enabled"`
	QueryOwnershipType  string   `graphql:"queryOwnership"      yaml:"queryOwnershipType"            json:"queryOwnershipType"`
	ThrottleTimeSeconds int      `graphql:"throttleTimeSeconds" yaml:"throttleTimeSeconds,omitempty" json:"throttleTimeSeconds,omitempty"`
	ThrottleField       string   `graphql:"throttleField"       yaml:"throttleField,omitempty"       json:"throttleField"`
	RunAsUserID         string   `graphql:"runAsUserId"         yaml:"runAsUserId,omitempty"         json:"runAsUserId,omitempty"`
}

type FilterAlerts struct {
	client *Client
}

func (c *Client) FilterAlerts() *FilterAlerts { return &FilterAlerts{client: c} }

func (fa *FilterAlerts) List(viewName string) ([]FilterAlert, error) {
	var query struct {
		SearchDomain struct {
			FilterAlerts []humiographql.FilterAlert `graphql:"filterAlerts"`
		} `graphql:"searchDomain(name: $viewName)"`
	}

	variables := map[string]any{
		"viewName": graphql.String(viewName),
	}

	err := fa.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	var filterAlerts = make([]FilterAlert, len(query.SearchDomain.FilterAlerts))
	for i := range query.SearchDomain.FilterAlerts {
		filterAlerts[i] = mapHumioGraphqlFilterAlertToFilterAlert(query.SearchDomain.FilterAlerts[i])
	}

	return filterAlerts, err
}

func (fa *FilterAlerts) Update(viewName string, updatedFilterAlert *FilterAlert) (*FilterAlert, error) {
	if updatedFilterAlert == nil {
		return nil, fmt.Errorf("updatedFilterAlert must not be nil")
	}

	if updatedFilterAlert.ID == "" {
		return nil, fmt.Errorf("updatedFilterAlert must have non-empty ID")
	}

	var mutation struct {
		humiographql.FilterAlert `graphql:"updateFilterAlert(input: $input)"`
	}

	actionNames := make([]graphql.String, len(updatedFilterAlert.ActionNames))
	for i, actionName := range updatedFilterAlert.ActionNames {
		actionNames[i] = graphql.String(actionName)
	}

	labels := make([]graphql.String, len(updatedFilterAlert.Labels))
	for i, label := range updatedFilterAlert.Labels {
		labels[i] = graphql.String(label)
	}

	updateAlert := humiographql.UpdateFilterAlert{
		ViewName:            humiographql.RepoOrViewName(viewName),
		ID:                  graphql.String(updatedFilterAlert.ID),
		Name:                graphql.String(updatedFilterAlert.Name),
		Description:         graphql.String(updatedFilterAlert.Description),
		QueryString:         graphql.String(updatedFilterAlert.QueryString),
		ActionIdsOrNames:    actionNames,
		Labels:              labels,
		Enabled:             graphql.Boolean(updatedFilterAlert.Enabled),
		RunAsUserID:         graphql.String(updatedFilterAlert.RunAsUserID),
		ThrottleTimeSeconds: humiographql.Long(updatedFilterAlert.ThrottleTimeSeconds),
		ThrottleField:       graphql.String(updatedFilterAlert.ThrottleField),
		QueryOwnershipType:  humiographql.QueryOwnershipType(updatedFilterAlert.QueryOwnershipType),
	}

	variables := map[string]any{
		"input": updateAlert,
	}

	err := fa.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	filterAlert := mapHumioGraphqlFilterAlertToFilterAlert(mutation.FilterAlert)

	return &filterAlert, nil
}

func (fa *FilterAlerts) Create(viewName string, newFilterAlert *FilterAlert) (*FilterAlert, error) {
	if newFilterAlert == nil {
		return nil, fmt.Errorf("newFilterAlert must not be nil")
	}

	var mutation struct {
		humiographql.FilterAlert `graphql:"createFilterAlert(input: $input)"`
	}

	actionNames := make([]graphql.String, len(newFilterAlert.ActionNames))
	for i, actionName := range newFilterAlert.ActionNames {
		actionNames[i] = graphql.String(actionName)
	}

	labels := make([]graphql.String, len(newFilterAlert.Labels))
	for i, label := range newFilterAlert.Labels {
		labels[i] = graphql.String(label)
	}

	createFilterAlert := humiographql.CreateFilterAlert{
		ViewName:            humiographql.RepoOrViewName(viewName),
		Name:                graphql.String(newFilterAlert.Name),
		Description:         graphql.String(newFilterAlert.Description),
		QueryString:         graphql.String(newFilterAlert.QueryString),
		ActionIdsOrNames:    actionNames,
		Labels:              labels,
		Enabled:             graphql.Boolean(newFilterAlert.Enabled),
		ThrottleTimeSeconds: humiographql.Long(newFilterAlert.ThrottleTimeSeconds),
		ThrottleField:       graphql.String(newFilterAlert.ThrottleField),
		RunAsUserID:         graphql.String(newFilterAlert.RunAsUserID),
		QueryOwnershipType:  humiographql.QueryOwnershipType(newFilterAlert.QueryOwnershipType),
	}

	variables := map[string]any{
		"input": createFilterAlert,
	}

	err := fa.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	filterAlert := mapHumioGraphqlFilterAlertToFilterAlert(mutation.FilterAlert)

	return &filterAlert, nil
}

func (fa *FilterAlerts) Delete(viewName, filterAlertID string) error {
	if filterAlertID == "" {
		return fmt.Errorf("filterAlertID is empty")
	}

	var mutation struct {
		DidDelete bool `graphql:"deleteFilterAlert(input: { viewName: $viewName, id: $id })"`
	}

	variables := map[string]any{
		"viewName": humiographql.RepoOrViewName(viewName),
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
			FilterAlert humiographql.FilterAlert `graphql:"filterAlert(id: $filterAlertId)"`
		} `graphql:"searchDomain(name: $viewName) "`
	}

	variables := map[string]any{
		"viewName":      graphql.String(viewName),
		"filterAlertId": graphql.String(filterAlertID),
	}

	err := fa.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	filterAlert := mapHumioGraphqlFilterAlertToFilterAlert(query.SearchDomain.FilterAlert)

	return &filterAlert, nil
}

func mapHumioGraphqlFilterAlertToFilterAlert(input humiographql.FilterAlert) FilterAlert {
	var queryOwnershipType, runAsUserID string
	switch input.QueryOwnership.QueryOwnershipTypeName {
	case humiographql.QueryOwnershipTypeNameOrganization:
		queryOwnershipType = QueryOwnershipTypeOrganization
	case humiographql.QueryOwnershipTypeNameUser:
		queryOwnershipType = QueryOwnershipTypeUser
		runAsUserID = string(input.QueryOwnership.ID)
	}

	var actionNames = make([]string, len(input.Actions))
	for i := range input.Actions {
		actionNames[i] = string(input.Actions[i].Name)
	}

	var labels = make([]string, len(input.Labels))
	for i := range input.Labels {
		labels[i] = string(input.Labels[i])
	}

	return FilterAlert{
		ID:                  string(input.ID),
		Name:                string(input.Name),
		Description:         string(input.Description),
		QueryString:         string(input.QueryString),
		ActionNames:         actionNames,
		Labels:              labels,
		Enabled:             bool(input.Enabled),
		ThrottleTimeSeconds: int(input.ThrottleTimeSeconds),
		ThrottleField:       string(input.ThrottleField),
		QueryOwnershipType:  queryOwnershipType,
		RunAsUserID:         runAsUserID,
	}
}
