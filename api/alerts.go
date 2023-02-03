package api

import (
	"fmt"
	"github.com/SaaldjorMike/graphql"
)

type Alert struct {
	ID                 string   `graphql:"id"                 yaml:"-"                     json:"id"`
	Name               string   `graphql:"name"               yaml:"name"                  json:"name"`
	QueryString        string   `graphql:"queryString"        yaml:"queryString"           json:"queryString"`
	QueryStart         string   `graphql:"queryStart"         yaml:"queryStart"            json:"queryStart"`
	ThrottleField      string   `graphql:"throttleField"      yaml:"throttleField"         json:"throttleField"`
	TimeOfLastTrigger  int      `graphql:"timeOfLastTrigger"  yaml:"timeOfLastTrigger"     json:"timeOfLastTrigger"`
	IsStarred          bool     `graphql:"isStarred"          yaml:"isStarred"             json:"isStarred"`
	Description        string   `graphql:"description"        yaml:"description,omitempty" json:"description"`
	ThrottleTimeMillis int      `graphql:"throttleTimeMillis" yaml:"throttleTimeMillis"    json:"throttleTimeMillis"`
	Enabled            bool     `graphql:"enabled"            yaml:"enabled"               json:"enabled"`
	Actions            []string `graphql:"actions"            yaml:"actions"               json:"actions"`
	Labels             []string `graphql:"labels"             yaml:"labels,omitempty"      json:"labels,omitempty"`
	LastError          string   `graphql:"lastError"          yaml:"lastError"             json:"lastError"`
}

type Long int64

type Alerts struct {
	client *Client
}

func (c *Client) Alerts() *Alerts { return &Alerts{client: c} }

func (a *Alerts) List(viewName string) ([]Alert, error) {
	var query struct {
		SearchDomain struct {
			Alerts []Alert `graphql:"alerts"`
		} `graphql:"searchDomain(name: $viewName)"`
	}

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
	}

	err := a.client.Query(&query, variables)
	return query.SearchDomain.Alerts, err
}

func (a *Alerts) Update(viewName string, newAlert *Alert) (*Alert, error) {
	if newAlert == nil {
		return nil, fmt.Errorf("newAlert must not be nil")
	}

	if newAlert.ID == "" {
		return nil, fmt.Errorf("newAlert must have non-empty newAlert id")
	}

	var mutation struct {
		Alert `graphql:"updateAlert(input: { id: $id, viewName: $viewName, name: $alertName, description: $description, queryString: $queryString, queryStart: $queryStart, throttleTimeMillis: $throttleTimeMillis, throttleField: $throttleField, enabled: $enabled, actions: $actions, labels: $labels })"`
	}

	actions := make([]graphql.String, len(newAlert.Actions))
	labels := make([]graphql.String, len(newAlert.Labels))
	var throttleField *graphql.String
	for i, action := range newAlert.Actions {
		actions[i] = graphql.String(action)
	}
	for i, label := range newAlert.Labels {
		labels[i] = graphql.String(label)
	}
	if newAlert.ThrottleField != "" {
		field := graphql.String(newAlert.ThrottleField)
		throttleField = &field
	}
	variables := map[string]interface{}{
		"id":                 graphql.String(newAlert.ID),
		"viewName":           graphql.String(viewName),
		"alertName":          graphql.String(newAlert.Name),
		"description":        graphql.String(newAlert.Description),
		"queryString":        graphql.String(newAlert.QueryString),
		"queryStart":         graphql.String(newAlert.QueryStart),
		"throttleField":      throttleField,
		"throttleTimeMillis": Long(newAlert.ThrottleTimeMillis),
		"enabled":            graphql.Boolean(newAlert.Enabled),
		"actions":            actions,
		"labels":             labels,
	}

	err := a.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	alert := Alert{
		ID:                 mutation.Alert.ID,
		Name:               mutation.Alert.Name,
		QueryString:        mutation.Alert.QueryString,
		QueryStart:         mutation.Alert.QueryStart,
		ThrottleField:      mutation.Alert.ThrottleField,
		TimeOfLastTrigger:  mutation.Alert.TimeOfLastTrigger,
		IsStarred:          mutation.Alert.IsStarred,
		Description:        mutation.Alert.Description,
		ThrottleTimeMillis: mutation.Alert.ThrottleTimeMillis,
		Enabled:            mutation.Alert.Enabled,
		Actions:            mutation.Alert.Actions,
		Labels:             mutation.Alert.Labels,
		LastError:          mutation.Alert.LastError,
	}

	return &alert, nil
}

func (a *Alerts) Add(viewName string, newAlert *Alert) (*Alert, error) {
	if newAlert == nil {
		return nil, fmt.Errorf("newAlert must not be nil")
	}
	var alert Alert

	var mutation struct {
		Alert `graphql:"createAlert(input: { viewName: $viewName, name: $alertName, description: $description, queryString: $queryString, queryStart: $queryStart, throttleTimeMillis: $throttleTimeMillis, throttleField: $throttleField, enabled: $enabled, actions: $actions, labels: $labels })"`
	}

	actions := make([]graphql.String, len(newAlert.Actions))
	labels := make([]graphql.String, len(newAlert.Labels))
	var throttleField *graphql.String
	for i, action := range newAlert.Actions {
		actions[i] = graphql.String(action)
	}
	for i, label := range newAlert.Labels {
		labels[i] = graphql.String(label)
	}
	if newAlert.ThrottleField != "" {
		field := graphql.String(newAlert.ThrottleField)
		throttleField = &field
	}
	variables := map[string]interface{}{
		"viewName":           graphql.String(viewName),
		"alertName":          graphql.String(newAlert.Name),
		"description":        graphql.String(newAlert.Description),
		"queryString":        graphql.String(newAlert.QueryString),
		"queryStart":         graphql.String(newAlert.QueryStart),
		"throttleTimeMillis": Long(newAlert.ThrottleTimeMillis),
		"throttleField":      throttleField,
		"enabled":            graphql.Boolean(newAlert.Enabled),
		"actions":            actions,
		"labels":             labels,
	}

	err := a.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	alert = Alert{
		ID:                 mutation.Alert.ID,
		Name:               mutation.Alert.Name,
		QueryString:        mutation.Alert.QueryString,
		QueryStart:         mutation.Alert.QueryStart,
		ThrottleField:      mutation.Alert.ThrottleField,
		TimeOfLastTrigger:  mutation.Alert.TimeOfLastTrigger,
		IsStarred:          mutation.Alert.IsStarred,
		Description:        mutation.Alert.Description,
		ThrottleTimeMillis: mutation.Alert.ThrottleTimeMillis,
		Enabled:            mutation.Alert.Enabled,
		Actions:            mutation.Alert.Actions,
		Labels:             mutation.Alert.Labels,
		LastError:          mutation.Alert.LastError,
	}
	return &alert, nil
}

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
