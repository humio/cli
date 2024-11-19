package api

import (
	"fmt"

	graphql "github.com/cli/shurcooL-graphql"
	"github.com/humio/cli/api/internal/humiographql"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type ScheduledSearch struct {
	ID                 string   `graphql:"id"                         yaml:"-"                          json:"id"`
	Name               string   `graphql:"name"                       yaml:"name"                       json:"name"`
	Description        string   `graphql:"description"                yaml:"description,omitempty"      json:"description,omitempty"`
	QueryString        string   `graphql:"queryString"                yaml:"queryString"                json:"queryString"`
	QueryStart         string   `graphql:"queryStart"                 yaml:"queryStart"                 json:"queryStart"`
	QueryEnd           string   `graphql:"queryEnd"                   yaml:"queryEnd"                   json:"queryEnd"`
	TimeZone           string   `graphql:"timeZone"                   yaml:"timeZone"                   json:"timeZone"`
	Schedule           string   `graphql:"schedule"                   yaml:"schedule"                   json:"schedule"`
	BackfillLimit      int      `graphql:"backfillLimit"              yaml:"backfillLimit"              json:"backfillLimit"`
	Enabled            bool     `graphql:"enabled"                    yaml:"enabled"                    json:"enabled"`
	ActionNames        []string `graphql:"actionNames"                yaml:"actionNames"                json:"actionNames"`
	RunAsUserID        string   `graphql:"runAsUserId"                yaml:"runAsUserId,omitempty"      json:"runAsUserId,omitempty"`
	Labels             []string `graphql:"labels"                     yaml:"labels"                     json:"labels"`
	QueryOwnershipType string   `graphql:"queryOwnership"             yaml:"queryOwnershipType"         json:"queryOwnershipType"`
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type ScheduledSearches struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) ScheduledSearches() *ScheduledSearches { return &ScheduledSearches{client: c} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *ScheduledSearches) List(viewName string) ([]ScheduledSearch, error) {
	var query struct {
		SearchDomain struct {
			ScheduledSearches []humiographql.ScheduledSearch `graphql:"scheduledSearches"`
		} `graphql:"searchDomain(name: $viewName)"`
	}

	variables := map[string]any{
		"viewName": graphql.String(viewName),
	}

	err := a.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	var scheduledSearches = make([]ScheduledSearch, len(query.SearchDomain.ScheduledSearches))
	for i := range query.SearchDomain.ScheduledSearches {
		scheduledSearches[i] = mapHumioGraphqlScheduledSearchToScheduledSearch(query.SearchDomain.ScheduledSearches[i])
	}

	return scheduledSearches, err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *ScheduledSearches) Update(viewName string, updateScheduledSearch *ScheduledSearch) (*ScheduledSearch, error) {
	if updateScheduledSearch == nil {
		return nil, fmt.Errorf("updateScheduledSearch must not be nil")
	}

	if updateScheduledSearch.ID == "" {
		return nil, fmt.Errorf("updateScheduledSearch must have non-empty ID")
	}

	var mutation struct {
		humiographql.ScheduledSearch `graphql:"updateScheduledSearch(input: $input)"`
	}

	actionNames := make([]graphql.String, len(updateScheduledSearch.ActionNames))
	for i, actionName := range updateScheduledSearch.ActionNames {
		actionNames[i] = graphql.String(actionName)
	}

	labels := make([]graphql.String, len(updateScheduledSearch.Labels))
	for i, label := range updateScheduledSearch.Labels {
		labels[i] = graphql.String(label)
	}

	updateAlert := humiographql.UpdateScheduledSearch{
		ViewName:          graphql.String(viewName),
		ID:                graphql.String(updateScheduledSearch.ID),
		Name:              graphql.String(updateScheduledSearch.Name),
		Description:       graphql.String(updateScheduledSearch.Description),
		QueryString:       graphql.String(updateScheduledSearch.QueryString),
		QueryStart:        graphql.String(updateScheduledSearch.QueryStart),
		QueryEnd:          graphql.String(updateScheduledSearch.QueryEnd),
		Schedule:          graphql.String(updateScheduledSearch.Schedule),
		TimeZone:          graphql.String(updateScheduledSearch.TimeZone),
		BackfillLimit:     graphql.Int(updateScheduledSearch.BackfillLimit),
		Enabled:           graphql.Boolean(updateScheduledSearch.Enabled),
		ActionsIdsOrNames: actionNames,
		Labels:            labels,
		RunAsUserID:       graphql.String(updateScheduledSearch.RunAsUserID),
		QueryOwnership:    humiographql.QueryOwnershipType(updateScheduledSearch.QueryOwnershipType),
	}

	variables := map[string]any{
		"input": updateAlert,
	}

	err := a.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	scheduledSearch := mapHumioGraphqlScheduledSearchToScheduledSearch(mutation.ScheduledSearch)

	return &scheduledSearch, nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *ScheduledSearches) Create(viewName string, newScheduledSearch *ScheduledSearch) (*ScheduledSearch, error) {
	if newScheduledSearch == nil {
		return nil, fmt.Errorf("newScheduledSearch must not be nil")
	}

	var mutation struct {
		humiographql.ScheduledSearch `graphql:"createScheduledSearch(input: $input)"`
	}

	actionNames := make([]graphql.String, len(newScheduledSearch.ActionNames))
	for i, actionName := range newScheduledSearch.ActionNames {
		actionNames[i] = graphql.String(actionName)
	}

	labels := make([]graphql.String, len(newScheduledSearch.Labels))
	for i, label := range newScheduledSearch.Labels {
		labels[i] = graphql.String(label)
	}

	createScheduledSearch := humiographql.CreateScheduledSearch{
		ViewName:          graphql.String(viewName),
		Name:              graphql.String(newScheduledSearch.Name),
		Description:       graphql.String(newScheduledSearch.Description),
		QueryString:       graphql.String(newScheduledSearch.QueryString),
		QueryStart:        graphql.String(newScheduledSearch.QueryStart),
		QueryEnd:          graphql.String(newScheduledSearch.QueryEnd),
		Schedule:          graphql.String(newScheduledSearch.Schedule),
		TimeZone:          graphql.String(newScheduledSearch.TimeZone),
		BackfillLimit:     graphql.Int(newScheduledSearch.BackfillLimit),
		Enabled:           graphql.Boolean(newScheduledSearch.Enabled),
		ActionsIdsOrNames: actionNames,
		Labels:            labels,
		RunAsUserID:       graphql.String(newScheduledSearch.RunAsUserID),
		QueryOwnership:    humiographql.QueryOwnershipType(newScheduledSearch.QueryOwnershipType),
	}

	variables := map[string]any{
		"input": createScheduledSearch,
	}

	err := a.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	scheduledSearch := mapHumioGraphqlScheduledSearchToScheduledSearch(mutation.ScheduledSearch)

	return &scheduledSearch, nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *ScheduledSearches) Delete(viewName, scheduledSearchID string) error {
	if scheduledSearchID == "" {
		return fmt.Errorf("scheduledSearchID is empty")
	}

	var mutation struct {
		DidDelete bool `graphql:"deleteScheduledSearch(input: { viewName: $viewName, id: $id })"`
	}

	variables := map[string]any{
		"viewName": graphql.String(viewName),
		"id":       graphql.String(scheduledSearchID),
	}

	err := a.client.Mutate(&mutation, variables)

	if !mutation.DidDelete {
		return fmt.Errorf("unable to remove scheduled search in repo/view '%s' with id '%s'", viewName, scheduledSearchID)
	}

	return err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (a *ScheduledSearches) Get(viewName string, scheduledSearchId string) (*ScheduledSearch, error) {
	var query struct {
		SearchDomain struct {
			ScheduledSearch humiographql.ScheduledSearch `graphql:"scheduledSearch(id: $scheduledSearchId)"`
		} `graphql:"searchDomain(name: $viewName) "`
	}

	variables := map[string]any{
		"viewName":          graphql.String(viewName),
		"scheduledSearchId": graphql.String(scheduledSearchId),
	}

	err := a.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	scheduledSearch := mapHumioGraphqlScheduledSearchToScheduledSearch(query.SearchDomain.ScheduledSearch)

	return &scheduledSearch, nil
}

func mapHumioGraphqlScheduledSearchToScheduledSearch(input humiographql.ScheduledSearch) ScheduledSearch {
	var queryOwnershipType, runAsUserID string
	switch input.QueryOwnership.QueryOwnershipTypeName {
	case humiographql.QueryOwnershipTypeNameOrganization:
		queryOwnershipType = QueryOwnershipTypeOrganization
	case humiographql.QueryOwnershipTypeNameUser:
		queryOwnershipType = QueryOwnershipTypeUser
		runAsUserID = string(input.QueryOwnership.ID)
	}

	var actionNames = make([]string, len(input.ActionsV2))
	for i := range input.ActionsV2 {
		actionNames[i] = string(input.ActionsV2[i].Name)
	}

	var labels = make([]string, len(input.Labels))
	for i := range input.Labels {
		labels[i] = string(input.Labels[i])
	}

	return ScheduledSearch{
		ID:                 string(input.ID),
		Name:               string(input.Name),
		Description:        string(input.Description),
		QueryString:        string(input.QueryString),
		QueryStart:         string(input.Start),
		QueryEnd:           string(input.End),
		TimeZone:           string(input.TimeZone),
		Schedule:           string(input.Schedule),
		BackfillLimit:      int(input.BackfillLimit),
		ActionNames:        actionNames,
		Labels:             labels,
		Enabled:            bool(input.Enabled),
		QueryOwnershipType: queryOwnershipType,
		RunAsUserID:        runAsUserID,
	}
}
