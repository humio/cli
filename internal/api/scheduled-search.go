package api

import (
	"context"
	"fmt"

	"github.com/humio/cli/internal/api/humiographql"
)

type ScheduledSearch struct {
	ID                 string `yaml:"-"`
	Name               string
	Description        *string
	QueryString        string `yaml:"queryString"`
	QueryStart         string `yaml:"queryStart"`
	QueryEnd           string `yaml:"queryEnd"`
	TimeZone           string `yaml:"timeZone"`
	Schedule           string
	BackfillLimit      int `yaml:"backfillLimit"`
	Enabled            bool
	ActionNames        []string `yaml:"actionNames"`
	OwnershipRunAsID   string   `yaml:"ownershipRunAsID"`
	Labels             []string
	QueryOwnershipType string `yaml:"queryOwnershipType"`
}

type ScheduledSearches struct {
	client *Client
}

func (c *Client) ScheduledSearches() *ScheduledSearches { return &ScheduledSearches{client: c} }

func (a *ScheduledSearches) List(searchDomainName string) ([]ScheduledSearch, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("searchDomainName must not be empty")
	}

	resp, err := humiographql.ListScheduledSearches(context.Background(), a.client, searchDomainName)
	if err != nil {
		return nil, err
	}
	respSearchDomain := resp.GetSearchDomain()
	respScheduledSearches := respSearchDomain.GetScheduledSearches()
	scheduledSearches := make([]ScheduledSearch, len(respScheduledSearches))
	for idx, scheduledSearch := range respScheduledSearches {
		actionNames := make([]string, len(scheduledSearch.GetActionsV2()))
		for kdx, action := range scheduledSearch.GetActionsV2() {
			actionNames[kdx] = action.GetName()
		}
		respQueryOwnership := scheduledSearch.GetQueryOwnership()
		runAsUserID := ""
		if respQueryOwnership != nil {
			runAsUserID = respQueryOwnership.GetId()
		}
		scheduledSearches[idx] = ScheduledSearch{
			ID:                 scheduledSearch.GetId(),
			Name:               scheduledSearch.GetName(),
			Description:        scheduledSearch.GetDescription(),
			QueryString:        scheduledSearch.GetQueryString(),
			QueryStart:         scheduledSearch.GetStart(),
			QueryEnd:           scheduledSearch.GetEnd(),
			Schedule:           scheduledSearch.GetSchedule(),
			TimeZone:           scheduledSearch.GetTimeZone(),
			BackfillLimit:      scheduledSearch.GetBackfillLimit(),
			Enabled:            scheduledSearch.GetEnabled(),
			ActionNames:        actionNames,
			Labels:             scheduledSearch.GetLabels(),
			OwnershipRunAsID:   runAsUserID,
			QueryOwnershipType: string(queryOwnershipToQueryOwnershipType(scheduledSearch.GetQueryOwnership())),
		}
	}
	return scheduledSearches, nil
}

func (a *ScheduledSearches) Update(searchDomainName string, updateScheduledSearch *ScheduledSearch) (*ScheduledSearch, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("searchDomainName must not be empty")
	}
	if updateScheduledSearch == nil {
		return nil, fmt.Errorf("updateScheduledSearch must not be nil")
	}

	if updateScheduledSearch.ID == "" {
		return nil, fmt.Errorf("updateScheduledSearch must have non-empty ID")
	}

	queryOwnershipType := humiographql.QueryOwnershipType(updateScheduledSearch.QueryOwnershipType)
	resp, err := humiographql.UpdateScheduledSearch(
		context.Background(),
		a.client,
		searchDomainName,
		updateScheduledSearch.ID,
		updateScheduledSearch.Name,
		updateScheduledSearch.Description,
		updateScheduledSearch.QueryString,
		updateScheduledSearch.QueryStart,
		updateScheduledSearch.QueryEnd,
		updateScheduledSearch.Schedule,
		updateScheduledSearch.TimeZone,
		updateScheduledSearch.BackfillLimit,
		updateScheduledSearch.Enabled,
		updateScheduledSearch.ActionNames,
		updateScheduledSearch.OwnershipRunAsID,
		updateScheduledSearch.Labels,
		&queryOwnershipType,
	)
	if err != nil {
		return nil, err
	}

	respScheduledSearch := resp.GetUpdateScheduledSearch()
	actionNames := make([]string, len(respScheduledSearch.GetActionsV2()))
	for kdx, action := range respScheduledSearch.GetActionsV2() {
		actionNames[kdx] = action.GetName()
	}
	respQueryOwnership := respScheduledSearch.GetQueryOwnership()
	runAsUserID := ""
	if respQueryOwnership != nil {
		runAsUserID = respQueryOwnership.GetId()
	}
	return &ScheduledSearch{
		ID:                 respScheduledSearch.GetId(),
		Name:               respScheduledSearch.GetName(),
		Description:        respScheduledSearch.GetDescription(),
		QueryString:        respScheduledSearch.GetQueryString(),
		QueryStart:         respScheduledSearch.GetStart(),
		QueryEnd:           respScheduledSearch.GetEnd(),
		TimeZone:           respScheduledSearch.GetTimeZone(),
		Schedule:           respScheduledSearch.GetSchedule(),
		BackfillLimit:      respScheduledSearch.GetBackfillLimit(),
		Enabled:            respScheduledSearch.GetEnabled(),
		ActionNames:        actionNames,
		OwnershipRunAsID:   runAsUserID,
		Labels:             respScheduledSearch.GetLabels(),
		QueryOwnershipType: string(queryOwnershipToQueryOwnershipType(respScheduledSearch.GetQueryOwnership())),
	}, nil
}

func (a *ScheduledSearches) Create(searchDomainName string, newScheduledSearch *ScheduledSearch) (*ScheduledSearch, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("searchDomainName must not be empty")
	}
	if newScheduledSearch == nil {
		return nil, fmt.Errorf("newScheduledSearch must not be nil")
	}

	var ownershipRunAsID *string
	if humiographql.QueryOwnershipType(newScheduledSearch.QueryOwnershipType) == humiographql.QueryOwnershipTypeUser {
		ownershipRunAsID = &newScheduledSearch.OwnershipRunAsID
	}

	queryOwnershipType := humiographql.QueryOwnershipType(newScheduledSearch.QueryOwnershipType)
	resp, err := humiographql.CreateScheduledSearch(
		context.Background(),
		a.client,
		searchDomainName,
		newScheduledSearch.Name,
		newScheduledSearch.Description,
		newScheduledSearch.QueryString,
		newScheduledSearch.QueryStart,
		newScheduledSearch.QueryEnd,
		newScheduledSearch.Schedule,
		newScheduledSearch.TimeZone,
		newScheduledSearch.BackfillLimit,
		newScheduledSearch.Enabled,
		newScheduledSearch.ActionNames,
		ownershipRunAsID,
		newScheduledSearch.Labels,
		&queryOwnershipType,
	)
	if err != nil {
		return nil, err
	}

	respScheduledSearch := resp.GetCreateScheduledSearch()
	actionNames := make([]string, len(respScheduledSearch.GetActionsV2()))
	for kdx, action := range respScheduledSearch.GetActionsV2() {
		actionNames[kdx] = action.GetName()
	}
	respQueryOwnership := respScheduledSearch.GetQueryOwnership()
	runAsUserID := ""
	if respQueryOwnership != nil {
		runAsUserID = respQueryOwnership.GetId()
	}
	return &ScheduledSearch{
		ID:                 respScheduledSearch.GetId(),
		Name:               respScheduledSearch.GetName(),
		Description:        respScheduledSearch.GetDescription(),
		QueryString:        respScheduledSearch.GetQueryString(),
		QueryStart:         respScheduledSearch.GetStart(),
		QueryEnd:           respScheduledSearch.GetEnd(),
		TimeZone:           respScheduledSearch.GetTimeZone(),
		Schedule:           respScheduledSearch.GetSchedule(),
		BackfillLimit:      respScheduledSearch.GetBackfillLimit(),
		Enabled:            respScheduledSearch.GetEnabled(),
		ActionNames:        actionNames,
		OwnershipRunAsID:   runAsUserID,
		Labels:             respScheduledSearch.GetLabels(),
		QueryOwnershipType: string(queryOwnershipToQueryOwnershipType(respScheduledSearch.GetQueryOwnership())),
	}, nil
}

func (a *ScheduledSearches) Delete(searchDomainName, scheduledSearchID string) error {
	if searchDomainName == "" {
		return fmt.Errorf("searchdomainName is empty")
	}
	if scheduledSearchID == "" {
		return fmt.Errorf("scheduledSearchID is empty")
	}

	_, err := humiographql.DeleteScheduledSearchByID(context.Background(), a.client, searchDomainName, scheduledSearchID)
	return err
}

func (a *ScheduledSearches) Get(searchDomainName string, scheduledSearchId string) (*ScheduledSearch, error) {
	resp, err := humiographql.GetScheduledSearchByID(context.Background(), a.client, searchDomainName, scheduledSearchId)
	if err != nil {
		return nil, err
	}

	respSearchDomain := resp.GetSearchDomain()
	respScheduledSearch := respSearchDomain.GetScheduledSearch()

	respActions := respScheduledSearch.GetActionsV2()
	actions := make([]string, len(respActions))
	for idx, action := range respActions {
		actions[idx] = action.GetName()
	}
	respQueryOwnership := respScheduledSearch.GetQueryOwnership()
	runAsUserID := ""
	if respQueryOwnership != nil {
		runAsUserID = respQueryOwnership.GetId()
	}
	return &ScheduledSearch{
		ID:                 respScheduledSearch.GetId(),
		Name:               respScheduledSearch.GetName(),
		Description:        respScheduledSearch.GetDescription(),
		QueryString:        respScheduledSearch.GetQueryString(),
		QueryStart:         respScheduledSearch.GetStart(),
		QueryEnd:           respScheduledSearch.GetEnd(),
		TimeZone:           respScheduledSearch.GetTimeZone(),
		Schedule:           respScheduledSearch.GetSchedule(),
		BackfillLimit:      respScheduledSearch.GetBackfillLimit(),
		Enabled:            respScheduledSearch.GetEnabled(),
		ActionNames:        actions,
		OwnershipRunAsID:   runAsUserID,
		Labels:             respScheduledSearch.GetLabels(),
		QueryOwnershipType: string(queryOwnershipToQueryOwnershipType(respScheduledSearch.GetQueryOwnership())),
	}, nil
}
