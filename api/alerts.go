package api

import (
	"fmt"

	graphql "github.com/cli/shurcooL-graphql"
	"github.com/humio/cli/api/internal/humiographql"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Alert struct {
	ID                 string   `graphql:"id"                 yaml:"-"                            json:"id"`
	Name               string   `graphql:"name"               yaml:"name"                         json:"name"`
	QueryString        string   `graphql:"queryString"        yaml:"queryString"                  json:"queryString"`
	QueryStart         string   `graphql:"queryStart"         yaml:"queryStart"                   json:"queryStart"`
	ThrottleField      string   `graphql:"throttleField"      yaml:"throttleField"                json:"throttleField"`
	TimeOfLastTrigger  int      `graphql:"timeOfLastTrigger"  yaml:"timeOfLastTrigger"            json:"timeOfLastTrigger"`
	IsStarred          bool     `graphql:"isStarred"          yaml:"isStarred"                    json:"isStarred"`
	Description        string   `graphql:"description"        yaml:"description,omitempty"        json:"description"`
	ThrottleTimeMillis int      `graphql:"throttleTimeMillis" yaml:"throttleTimeMillis"           json:"throttleTimeMillis"`
	Enabled            bool     `graphql:"enabled"            yaml:"enabled"                      json:"enabled"`
	Actions            []string `graphql:"actions"            yaml:"actions"                      json:"actions"`
	Labels             []string `graphql:"labels"             yaml:"labels,omitempty"             json:"labels,omitempty"`
	LastError          string   `graphql:"lastError"          yaml:"lastError"                    json:"lastError"`
	RunAsUserID        string   `graphql:"runAsUserId"        yaml:"runAsUserId,omitempty"        json:"runAsUserId,omitempty"`
	QueryOwnershipType string   `graphql:"queryOwnershipType" yaml:"queryOwnershipType,omitempty" json:"queryOwnershipType,omitempty"`
}

const (
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	QueryOwnershipTypeUser string = "User"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	QueryOwnershipTypeOrganization string = "Organization"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Alerts struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Alerts() *Alerts { return &Alerts{client: c} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *Alerts) List(viewName string) ([]Alert, error) {
	var query struct {
		SearchDomain struct {
			Alerts []humiographql.Alert `graphql:"alerts"`
		} `graphql:"searchDomain(name: $viewName)"`
	}

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
	}

	err := a.client.Query(&query, variables)

	var alerts []Alert
	for _, humioGraphqlAlert := range query.SearchDomain.Alerts {
		alerts = append(alerts, mapHumioGraphqlAlertToAlert(humioGraphqlAlert))
	}
	return alerts, err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *Alerts) Update(viewName string, newAlert *Alert) (*Alert, error) {
	if newAlert == nil {
		return nil, fmt.Errorf("newAlert must not be nil")
	}

	if newAlert.ID == "" {
		return nil, fmt.Errorf("newAlert must have non-empty newAlert id")
	}

	var mutation struct {
		humiographql.Alert `graphql:"updateAlert(input: $input)"`
	}

	actions := make([]graphql.String, len(newAlert.Actions))
	for i, action := range newAlert.Actions {
		actions[i] = graphql.String(action)
	}

	labels := make([]graphql.String, len(newAlert.Labels))
	for i, label := range newAlert.Labels {
		labels[i] = graphql.String(label)
	}

	updateAlert := humiographql.UpdateAlert{
		ID:                 graphql.String(newAlert.ID),
		ViewName:           graphql.String(viewName),
		Name:               graphql.String(newAlert.Name),
		Description:        graphql.String(newAlert.Description),
		QueryString:        graphql.String(newAlert.QueryString),
		QueryStart:         graphql.String(newAlert.QueryStart),
		ThrottleTimeMillis: humiographql.Long(newAlert.ThrottleTimeMillis),
		Enabled:            graphql.Boolean(newAlert.Enabled),
		Actions:            actions,
		Labels:             labels,
		RunAsUserID:        graphql.String(newAlert.RunAsUserID),
		QueryOwnershipType: humiographql.QueryOwnershipType(newAlert.QueryOwnershipType),
		ThrottleField:      graphql.String(newAlert.ThrottleField),
	}

	variables := map[string]interface{}{
		"input": updateAlert,
	}

	err := a.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	alert := mapHumioGraphqlAlertToAlert(mutation.Alert)

	return &alert, nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *Alerts) Add(viewName string, newAlert *Alert) (*Alert, error) {
	if newAlert == nil {
		return nil, fmt.Errorf("newAlert must not be nil")
	}

	var mutation struct {
		humiographql.Alert `graphql:"createAlert(input: $input)"`
	}

	actions := make([]graphql.String, len(newAlert.Actions))
	for i, action := range newAlert.Actions {
		actions[i] = graphql.String(action)
	}

	labels := make([]graphql.String, len(newAlert.Labels))
	for i, label := range newAlert.Labels {
		labels[i] = graphql.String(label)
	}

	createAlert := humiographql.CreateAlert{
		ViewName:           graphql.String(viewName),
		Name:               graphql.String(newAlert.Name),
		Description:        graphql.String(newAlert.Description),
		QueryString:        graphql.String(newAlert.QueryString),
		QueryStart:         graphql.String(newAlert.QueryStart),
		ThrottleTimeMillis: humiographql.Long(newAlert.ThrottleTimeMillis),
		Enabled:            graphql.Boolean(newAlert.Enabled),
		Actions:            actions,
		Labels:             labels,
		RunAsUserID:        graphql.String(newAlert.RunAsUserID),
		QueryOwnershipType: humiographql.QueryOwnershipType(newAlert.QueryOwnershipType),
		ThrottleField:      graphql.String(newAlert.ThrottleField),
	}

	variables := map[string]interface{}{
		"input": createAlert,
	}

	err := a.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	alert := mapHumioGraphqlAlertToAlert(mutation.Alert)
	return &alert, nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *Alerts) Get(viewName, alertName string) (*Alert, error) {
	alerts, err := a.List(viewName)
	if err != nil {
		return nil, fmt.Errorf("unable to list alerts: %w", err)
	}
	for _, alert := range alerts {
		if alert.Name == alertName {
			return &alert, nil
		}
	}

	return nil, AlertNotFound(alertName)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *Alerts) Delete(viewName, alertName string) error {
	actions, err := a.List(viewName)
	if err != nil {
		return fmt.Errorf("unable to list alerts: %w", err)
	}
	var alertId string
	for _, alert := range actions {
		if alert.Name == alertName {
			alertId = alert.ID
			break
		}
	}
	if alertId == "" {
		return fmt.Errorf("unable to find alert")
	}

	var mutation struct {
		DeleteAction bool `graphql:"deleteAlert(input: { viewName: $viewName, id: $alertId })"`
	}

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
		"alertId":  graphql.String(alertId),
	}

	return a.client.Mutate(&mutation, variables)
}

func mapHumioGraphqlAlertToAlert(input humiographql.Alert) Alert {
	var queryOwnershipType string
	switch input.QueryOwnership.QueryOwnershipTypeName {
	case humiographql.QueryOwnershipTypeNameOrganization:
		queryOwnershipType = QueryOwnershipTypeOrganization
	case humiographql.QueryOwnershipTypeNameUser:
		queryOwnershipType = QueryOwnershipTypeUser
	}

	var actions []string
	for _, action := range input.Actions {
		actions = append(actions, string(action))
	}

	var labels []string
	for _, label := range input.Labels {
		labels = append(labels, string(label))
	}

	return Alert{
		ID:                 string(input.ID),
		Name:               string(input.Name),
		QueryString:        string(input.QueryString),
		QueryStart:         string(input.QueryStart),
		ThrottleField:      string(input.ThrottleField),
		TimeOfLastTrigger:  int(input.TimeOfLastTrigger),
		IsStarred:          bool(input.IsStarred),
		Description:        string(input.Description),
		ThrottleTimeMillis: int(input.ThrottleTimeMillis),
		Enabled:            bool(input.Enabled),
		Actions:            actions,
		Labels:             labels,
		LastError:          string(input.LastError),
		RunAsUserID:        string(input.RunAsUser.ID),
		QueryOwnershipType: queryOwnershipType,
	}
}
