package api

import (
	"fmt"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/humio/cli/api/internal/humiographql"
)

type FilterAlert struct {
	ID                 string   `graphql:"id"             yaml:"id"                      json:"id"`
	Name               string   `graphql:"name"           yaml:"name"                    json:"name"`
	Description        string   `graphql:"description"    yaml:"description,omitempty"   json:"description,omitempty"`
	QueryString        string   `graphql:"queryString"    yaml:"queryString"             json:"queryString"`
	Actions            []string `graphql:"actions"        yaml:"actions"                 json:"actions"`
	Labels             []string `graphql:"labels"         yaml:"labels"                  json:"labels"`
	Enabled            bool     `graphql:"enabled"        yaml:"enabled"                 json:"enabled"`
	LastTriggered      int      `graphql:"lastTriggered"  yaml:"lastTriggered,omitempty" json:"lastTriggered,omitempty"`
	LastErrorTime      int      `graphql:"lastErrorTime"  yaml:"lastErrorTime,omitempty" json:"lastErrorTime,omitempty"`
	LastError          string   `graphql:"lastError"      yaml:"lastError,omitempty"     json:"lastError,omitempty"`
	LastWarnings       []string `graphql:"lastWarnings"   yaml:"lastWarnings"            json:"lastWarnings"`
	QueryOwnershipType string   `graphql:"queryOwnership" yaml:"queryOwnershipType"      json:"queryOwnershipType"`
	RunAsUserID        string   `graphql:"runAsUserID"    yaml:"runAsUserID,omitempty"   json:"runAsUserID,omitempty"`
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

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
	}

	err := fa.client.Query(&query, variables)

	var filterAlerts []FilterAlert
	for _, humioGraphqlFilterAlert := range query.SearchDomain.FilterAlerts {
		filterAlerts = append(filterAlerts, mapHumioGraphqlFilterAlertToFilterAlert(humioGraphqlFilterAlert))
	}
	return filterAlerts, err
}

func (fa *FilterAlerts) Update(viewName string, newFilterAlert *FilterAlert) (*FilterAlert, error) {
	if newFilterAlert == nil {
		return nil, fmt.Errorf("newFilterAlert must not be nil")
	}

	if newFilterAlert.ID == "" {
		return nil, fmt.Errorf("newFilterAlert must have non-empty newFilterAlert id")
	}

	var mutation struct {
		humiographql.FilterAlert `graphql:"updateFilterAlert(input: $input)"`
	}

	actionsIdsOrNames := make([]graphql.String, len(newFilterAlert.Actions))
	for i, actionIdOrName := range newFilterAlert.Actions {
		actionsIdsOrNames[i] = graphql.String(actionIdOrName)
	}

	labels := make([]graphql.String, len(newFilterAlert.Labels))
	for i, label := range newFilterAlert.Labels {
		labels[i] = graphql.String(label)
	}

	updateAlert := humiographql.UpdateFilterAlert{
		ViewName:           humiographql.RepoOrViewName(viewName),
		ID:                 graphql.String(newFilterAlert.ID),
		Name:               graphql.String(newFilterAlert.Name),
		Description:        graphql.String(newFilterAlert.Description),
		QueryString:        graphql.String(newFilterAlert.QueryString),
		ActionIdsOrNames:   actionsIdsOrNames,
		Labels:             labels,
		Enabled:            graphql.Boolean(newFilterAlert.Enabled),
		RunAsUserID:        graphql.String(newFilterAlert.RunAsUserID),
		QueryOwnershipType: humiographql.QueryOwnershipType(newFilterAlert.QueryOwnershipType),
	}

	variables := map[string]interface{}{
		"input": updateAlert,
	}

	err := fa.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	filterAlert := mapHumioGraphqlFilterAlertToFilterAlert(mutation.FilterAlert)

	return &filterAlert, nil
}

func (fa *FilterAlerts) Add(viewName string, newFilterAlert *FilterAlert) (*FilterAlert, error) {
	if newFilterAlert == nil {
		return nil, fmt.Errorf("newFilterAlert must not be nil")
	}

	var mutation struct {
		humiographql.FilterAlert `graphql:"createFilterAlert(input: $input)"`
	}

	actionsIdsOrNames := make([]graphql.String, len(newFilterAlert.Actions))
	for i, actionIdOrName := range newFilterAlert.Actions {
		actionsIdsOrNames[i] = graphql.String(actionIdOrName)
	}

	labels := make([]graphql.String, len(newFilterAlert.Labels))
	for i, label := range newFilterAlert.Labels {
		labels[i] = graphql.String(label)
	}

	createFilterAlert := humiographql.CreateFilterAlert{
		ViewName:           humiographql.RepoOrViewName(viewName),
		Name:               graphql.String(newFilterAlert.Name),
		Description:        graphql.String(newFilterAlert.Description),
		QueryString:        graphql.String(newFilterAlert.QueryString),
		ActionIdsOrNames:   actionsIdsOrNames,
		Labels:             labels,
		Enabled:            graphql.Boolean(newFilterAlert.Enabled),
		RunAsUserID:        graphql.String(newFilterAlert.RunAsUserID),
		QueryOwnershipType: humiographql.QueryOwnershipType(newFilterAlert.QueryOwnershipType),
	}

	variables := map[string]interface{}{
		"input": createFilterAlert,
	}

	err := fa.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	filterAlert := mapHumioGraphqlFilterAlertToFilterAlert(mutation.FilterAlert)

	return &filterAlert, nil
}

func (fa *FilterAlerts) Get(viewName, filterAlertName string) (*FilterAlert, error) {
	filterAlerts, err := fa.List(viewName)
	if err != nil {
		return nil, fmt.Errorf("unable to list filter alerts: %w", err)
	}
	for _, filterAlert := range filterAlerts {
		if filterAlert.Name == filterAlertName {
			return &filterAlert, nil
		}
	}

	return nil, FilterAlertNotFound(filterAlertName)
}

func (fa *FilterAlerts) Delete(viewName, filterAlertName string) error {
	filterAlerts, err := fa.List(viewName)
	if err != nil {
		return fmt.Errorf("unable to list filter alerts: %w", err)
	}
	var filterAlertID string
	for _, filterAlert := range filterAlerts {
		if filterAlert.Name == filterAlertName {
			filterAlertID = filterAlert.ID
			break
		}
	}
	if filterAlertID == "" {
		return fmt.Errorf("unable to find filter alert")
	}

	var mutation struct {
		DeleteAction bool `graphql:"deleteFilterAlert(input: { viewName: $viewName, id: $id })"`
	}

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
		"id":       graphql.String(filterAlertID),
	}

	return fa.client.Mutate(&mutation, variables)
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

	var actions []string
	for _, action := range input.Actions {
		actions = append(actions, string(action.ID))
	}

	var labels []string
	for _, label := range input.Labels {
		labels = append(labels, string(label))
	}

	var lastWarnings []string
	for _, warning := range input.LastWarnings {
		lastWarnings = append(lastWarnings, string(warning))
	}

	return FilterAlert{
		ID:                 string(input.ID),
		Name:               string(input.Name),
		Description:        string(input.Description),
		QueryString:        string(input.QueryString),
		Actions:            actions,
		Labels:             labels,
		Enabled:            bool(input.Enabled),
		LastTriggered:      int(input.LastTriggered),
		LastErrorTime:      int(input.LastErrorTime),
		LastError:          string(input.LastError),
		LastWarnings:       lastWarnings,
		QueryOwnershipType: queryOwnershipType,
		RunAsUserID:        runAsUserID,
	}
}
