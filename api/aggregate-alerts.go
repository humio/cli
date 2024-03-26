package api

import (
	"fmt"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/humio/cli/api/internal/humiographql"
)

type AggregateAlert struct {
	ID                    string   `graphql:"id"                    yaml:"-"                       json:"id"`
	Name                  string   `graphql:"name"                  yaml:"name"                    json:"name"`
	Description           string   `graphql:"description"           yaml:"description,omitempty"   json:"description,omitempty"`
	QueryString           string   `graphql:"queryString"           yaml:"queryString"             json:"queryString"`
	SearchIntervalSeconds int      `graphql:"searchIntervalSeconds" yaml:"searchIntervalSeconds"   json:"searchIntervalSeconds"`
	ActionNames           []string `graphql:"actionNames"           yaml:"actionNames"             json:"actionNames"`
	Labels                []string `graphql:"labels"                yaml:"labels"                  json:"labels"`
	Enabled               bool     `graphql:"enabled"               yaml:"enabled"                 json:"enabled"`
	ThrottleField         string   `graphql:"throttleField"         yaml:"throttleField,omitempty" json:"throttleField,omitempty"`
	ThrottleTimeSeconds   int      `graphql:"throttleTimeSeconds"   yaml:"throttleTimeSeconds"     json:"throttleTimeSeconds"`
	QueryOwnershipType    string   `graphql:"queryOwnership"        yaml:"queryOwnershipType"      json:"queryOwnershipType"`
	TriggerMode           string   `graphql:"triggerMode"           yaml:"triggerMode"             json:"triggerMode"`
	QueryTimestampType    string   `graphql:"queryTimestampType"    yaml:"queryTimestampType"      json:"queryTimestampType"`
	RunAsUserID           string   `graphql:"runAsUserId"           yaml:"runAsUserId,omitempty"   json:"runAsUserId,omitempty"`
}

type AggregateAlerts struct {
	client *Client
}

func (c *Client) AggregateAlerts() *AggregateAlerts { return &AggregateAlerts{client: c} }

func (a *AggregateAlerts) List(viewName string) ([]AggregateAlert, error) {
	var query struct {
		SearchDomain struct {
			AggregateAlerts []humiographql.AggregateAlert `graphql:"aggregateAlerts"`
		} `graphql:"searchDomain(name: $viewName)"`
	}

	variables := map[string]any{
		"viewName": graphql.String(viewName),
	}

	err := a.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	var aggregateAlerts = make([]AggregateAlert, len(query.SearchDomain.AggregateAlerts))
	for i := range query.SearchDomain.AggregateAlerts {
		aggregateAlerts[i] = mapHumioGraphqlAggregateAlertToAggregateAlert(query.SearchDomain.AggregateAlerts[i])
	}

	return aggregateAlerts, err
}

func (a *AggregateAlerts) Update(viewName string, updatedAggregateAlert *AggregateAlert) (*AggregateAlert, error) {
	if updatedAggregateAlert == nil {
		return nil, fmt.Errorf("updatedAggregateAlert must not be nil")
	}

	if updatedAggregateAlert.ID == "" {
		return nil, fmt.Errorf("updatedAggregateAlert must have non-empty ID")
	}

	var mutation struct {
		humiographql.AggregateAlert `graphql:"updateAggregateAlert(input: $input)"`
	}

	actionNames := make([]graphql.String, len(updatedAggregateAlert.ActionNames))
	for i, actionName := range updatedAggregateAlert.ActionNames {
		actionNames[i] = graphql.String(actionName)
	}

	labels := make([]graphql.String, len(updatedAggregateAlert.Labels))
	for i, label := range updatedAggregateAlert.Labels {
		labels[i] = graphql.String(label)
	}

	updateAlert := humiographql.UpdateAggregateAlert{
		ViewName:              humiographql.RepoOrViewName(viewName),
		ID:                    graphql.String(updatedAggregateAlert.ID),
		Name:                  graphql.String(updatedAggregateAlert.Name),
		Description:           graphql.String(updatedAggregateAlert.Description),
		QueryString:           graphql.String(updatedAggregateAlert.QueryString),
		SearchIntervalSeconds: humiographql.Long(updatedAggregateAlert.SearchIntervalSeconds),
		ActionIdsOrNames:      actionNames,
		Labels:                labels,
		Enabled:               graphql.Boolean(updatedAggregateAlert.Enabled),
		RunAsUserID:           graphql.String(updatedAggregateAlert.RunAsUserID),
		ThrottleField:         graphql.String(updatedAggregateAlert.ThrottleField),
		ThrottleTimeSeconds:   humiographql.Long(updatedAggregateAlert.ThrottleTimeSeconds),
		TriggerMode:           humiographql.TriggerMode(updatedAggregateAlert.TriggerMode),
		QueryTimestampType:    humiographql.QueryTimestampType(updatedAggregateAlert.QueryTimestampType),
		QueryOwnershipType:    humiographql.QueryOwnershipType(updatedAggregateAlert.QueryOwnershipType),
	}

	variables := map[string]any{
		"input": updateAlert,
	}

	err := a.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	aggregateAlert := mapHumioGraphqlAggregateAlertToAggregateAlert(mutation.AggregateAlert)

	return &aggregateAlert, nil
}

func (a *AggregateAlerts) Create(viewName string, newAggregateAlert *AggregateAlert) (*AggregateAlert, error) {
	if newAggregateAlert == nil {
		return nil, fmt.Errorf("newAggregateAlert must not be nil")
	}

	var mutation struct {
		humiographql.AggregateAlert `graphql:"createAggregateAlert(input: $input)"`
	}

	actionNames := make([]graphql.String, len(newAggregateAlert.ActionNames))
	for i, actionName := range newAggregateAlert.ActionNames {
		actionNames[i] = graphql.String(actionName)
	}

	labels := make([]graphql.String, len(newAggregateAlert.Labels))
	for i, label := range newAggregateAlert.Labels {
		labels[i] = graphql.String(label)
	}

	createAggregateAlert := humiographql.CreateAggregateAlert{
		ViewName:              humiographql.RepoOrViewName(viewName),
		Name:                  graphql.String(newAggregateAlert.Name),
		Description:           graphql.String(newAggregateAlert.Description),
		QueryString:           graphql.String(newAggregateAlert.QueryString),
		SearchIntervalSeconds: humiographql.Long(newAggregateAlert.SearchIntervalSeconds),
		ActionIdsOrNames:      actionNames,
		Labels:                labels,
		Enabled:               graphql.Boolean(newAggregateAlert.Enabled),
		ThrottleField:         graphql.String(newAggregateAlert.ThrottleField),
		ThrottleTimeSeconds:   humiographql.Long(newAggregateAlert.ThrottleTimeSeconds),
		RunAsUserID:           graphql.String(newAggregateAlert.RunAsUserID),
		TriggerMode:           humiographql.TriggerMode(newAggregateAlert.TriggerMode),
		QueryTimestampType:    humiographql.QueryTimestampType(newAggregateAlert.QueryTimestampType),
		QueryOwnershipType:    humiographql.QueryOwnershipType(newAggregateAlert.QueryOwnershipType),
	}

	variables := map[string]any{
		"input": createAggregateAlert,
	}

	err := a.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	aggregateAlert := mapHumioGraphqlAggregateAlertToAggregateAlert(mutation.AggregateAlert)

	return &aggregateAlert, nil
}

func (a *AggregateAlerts) Delete(viewName, aggregateAlertID string) error {
	if aggregateAlertID == "" {
		return fmt.Errorf("aggregateAlertID is empty")
	}

	var mutation struct {
		DidDelete bool `graphql:"deleteAggregateAlert(input: { viewName: $viewName, id: $id })"`
	}

	variables := map[string]any{
		"viewName": humiographql.RepoOrViewName(viewName),
		"id":       graphql.String(aggregateAlertID),
	}

	err := a.client.Mutate(&mutation, variables)

	if !mutation.DidDelete {
		return fmt.Errorf("unable to remove aggregate alert in repo/view '%s' with id '%s'", viewName, aggregateAlertID)
	}

	return err
}

func (a *AggregateAlerts) Get(viewName string, aggregateAlertID string) (*AggregateAlert, error) {
	var query struct {
		SearchDomain struct {
			AggregateAlert humiographql.AggregateAlert `graphql:"aggregateAlert(id: $aggregateAlertId)"`
		} `graphql:"searchDomain(name: $viewName) "`
	}

	variables := map[string]any{
		"viewName":         graphql.String(viewName),
		"aggregateAlertId": graphql.String(aggregateAlertID),
	}

	err := a.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	aggregateAlert := mapHumioGraphqlAggregateAlertToAggregateAlert(query.SearchDomain.AggregateAlert)

	return &aggregateAlert, nil
}

func mapHumioGraphqlAggregateAlertToAggregateAlert(input humiographql.AggregateAlert) AggregateAlert {
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

	return AggregateAlert{
		ID:                    string(input.ID),
		Name:                  string(input.Name),
		Description:           string(input.Description),
		QueryString:           string(input.QueryString),
		SearchIntervalSeconds: int(input.SearchIntervalSeconds),
		ActionNames:           actionNames,
		Labels:                labels,
		Enabled:               bool(input.Enabled),
		ThrottleField:         string(input.ThrottleField),
		ThrottleTimeSeconds:   int(input.ThrottleTimeSeconds),
		QueryOwnershipType:    queryOwnershipType,
		TriggerMode:           string(input.TriggerMode),
		QueryTimestampType:    string(input.QueryTimestampType),
		RunAsUserID:           runAsUserID,
	}
}
