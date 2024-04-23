package api

import (
	"fmt"

	graphql "github.com/cli/shurcooL-graphql"
)

type Alert struct {
	ID                 string         `graphql:"id"                 yaml:"-"                            json:"id"`
	Name               string         `graphql:"name"               yaml:"name"                         json:"name"`
	DisplayName        string         `graphql:"displayName" yaml:"displayName" json:"displayName"`
	QueryString        string         `graphql:"queryString"        yaml:"queryString"                  json:"queryString"`
	QueryStart         string         `graphql:"queryStart"         yaml:"queryStart"                   json:"queryStart"`
	ThrottleField      string         `graphql:"throttleField"      yaml:"throttleField"                json:"throttleField"`
	TimeOfLastTrigger  int64          `graphql:"timeOfLastTrigger"  yaml:"timeOfLastTrigger"            json:"timeOfLastTrigger"`
	IsStarred          bool           `graphql:"isStarred"          yaml:"isStarred"                    json:"isStarred"`
	Description        string         `graphql:"description"        yaml:"description,omitempty"        json:"description"`
	ThrottleTimeMillis int64          `graphql:"throttleTimeMillis" yaml:"throttleTimeMillis"           json:"throttleTimeMillis"`
	Enabled            bool           `graphql:"enabled"            yaml:"enabled"                      json:"enabled"`
	Actions            []string       `graphql:"actions"            yaml:"actions"                      json:"actions"`
	Labels             []string       `graphql:"labels"             yaml:"labels,omitempty"             json:"labels,omitempty"`
	LastError          string         `graphql:"lastError"          yaml:"lastError"                    json:"lastError"`
	RunAsUser          User           `graphql:"runAsUser"          yaml:"runAsUser,omitempty"          json:"runAsUser,omitempty"`
	QueryOwnership     QueryOwnership `graphql:"queryOwnership"     yaml:"queryOwnership"               json:"queryOwnership"`
}

type QueryOwnership struct {
	ID                    string `graphql:"id"`
	UserOwnership         `graphql:"... on UserOwnership"`
	OrganizationOwnership `graphql:"... on OrganizationOwnership"`
}

type UserOwnership struct {
	User User `graphql:"user"`
}

type OrganizationOwnership struct {
	Organization Organization `graphql:"organization"`
}

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

func (a *Alerts) Update(viewName string, newAlert UpdateAlert) (*Alert, error) {
	var mutation struct {
		Alert Alert `graphql:"updateAlert(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": newAlert,
	}

	err := a.client.Mutate(&mutation, variables)
	return &mutation.Alert, err
}

func (a *Alerts) Add(viewName string, newAlert CreateAlert) (*Alert, error) {

	var mutation struct {
		Alert Alert `graphql:"createAlert(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": newAlert,
	}

	err := a.client.Mutate(&mutation, variables)
	return &mutation.Alert, err
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
