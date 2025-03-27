package api

import (
	"context"
	"fmt"

	"github.com/humio/cli/internal/api/humiographql"
)

type ScheduledSearchV2 struct {
	ID                          string   `yaml:"-"`
	Name                        string   `yaml:"name"`
	Description                 *string  `yaml:"description,omitempty"`
	QueryString                 string   `yaml:"queryString"`
	SearchIntervalSeconds       int64    `yaml:"searchIntervalSeconds"`
	SearchIntervalOffsetSeconds *int64   `yaml:"searchIntervalOffsetSeconds,omitempty"`
	QueryTimestampType          string   `yaml:"queryTimestampType"`
	MaxWaitTimeSeconds          *int64   `yaml:"maxWaitTimeSeconds,omitempty"`
	BackfillLimitV2             *int     `yaml:"backfillLimitV2,omitempty"`
	TimeZone                    string   `yaml:"timeZone"`
	Schedule                    string   `yaml:"schedule"`
	Enabled                     bool     `yaml:"enabled"`
	ActionNames                 []string `yaml:"actionNames"`
	Labels                      []string `yaml:"labels"`
	QueryOwnershipType          string   `yaml:"queryOwnershipType"`
	OwnershipRunAsID            string   `yaml:"ownershipRunAsID"`
}

type ScheduledSearchesV2 struct {
	client *Client
}

func (c *Client) ScheduledSearchesV2() *ScheduledSearchesV2 { return &ScheduledSearchesV2{client: c} }

func (a *ScheduledSearchesV2) List(searchDomainName string) ([]ScheduledSearchV2, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("searchDomainName must not be empty")
	}

	resp, err := humiographql.ListScheduledSearchesV2(context.Background(), a.client, searchDomainName)
	if err != nil {
		return nil, err
	}
	respSearchDomain := resp.GetSearchDomain()
	respScheduledSearches := respSearchDomain.GetScheduledSearches()
	scheduledSearches := make([]ScheduledSearchV2, len(respScheduledSearches))
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
		scheduledSearches[idx] = ScheduledSearchV2{
			ID:                          scheduledSearch.GetId(),
			Name:                        scheduledSearch.GetName(),
			Description:                 scheduledSearch.GetDescription(),
			QueryString:                 scheduledSearch.GetQueryString(),
			SearchIntervalSeconds:       scheduledSearch.GetSearchIntervalSeconds(),
			SearchIntervalOffsetSeconds: scheduledSearch.GetSearchIntervalOffsetSeconds(),
			MaxWaitTimeSeconds:          scheduledSearch.GetMaxWaitTimeSeconds(),
			Schedule:                    scheduledSearch.GetSchedule(),
			TimeZone:                    scheduledSearch.GetTimeZone(),
			BackfillLimitV2:             scheduledSearch.GetBackfillLimitV2(),
			Enabled:                     scheduledSearch.GetEnabled(),
			ActionNames:                 actionNames,
			Labels:                      scheduledSearch.GetLabels(),
			OwnershipRunAsID:            runAsUserID,
			QueryTimestampType:          string(scheduledSearch.QueryTimestampType),
			QueryOwnershipType:          string(queryOwnershipToQueryOwnershipType(scheduledSearch.GetQueryOwnership())),
		}
	}
	return scheduledSearches, nil
}

func (a *ScheduledSearchesV2) Create(searchDomainName string, newScheduledSearch *ScheduledSearchV2) (*ScheduledSearchV2, error) {
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
	queryTimestampType := humiographql.QueryTimestampType(newScheduledSearch.QueryTimestampType)

	resp, err := humiographql.CreateScheduledSearchV2(
		context.Background(),
		a.client,
		searchDomainName,
		newScheduledSearch.Name,
		newScheduledSearch.Description,
		newScheduledSearch.QueryString,
		newScheduledSearch.SearchIntervalSeconds,
		newScheduledSearch.SearchIntervalOffsetSeconds,
		newScheduledSearch.MaxWaitTimeSeconds,
		newScheduledSearch.Schedule,
		newScheduledSearch.TimeZone,
		newScheduledSearch.BackfillLimitV2,
		newScheduledSearch.Enabled,
		newScheduledSearch.ActionNames,
		ownershipRunAsID,
		newScheduledSearch.Labels,
		queryTimestampType,
		queryOwnershipType,
	)
	if err != nil {
		return nil, err
	}

	respScheduledSearch := resp.GetCreateScheduledSearchV2()
	actionNames := make([]string, len(respScheduledSearch.GetActionsV2()))
	for kdx, action := range respScheduledSearch.GetActionsV2() {
		actionNames[kdx] = action.GetName()
	}
	respQueryOwnership := respScheduledSearch.GetQueryOwnership()
	runAsUserID := ""
	if respQueryOwnership != nil {
		runAsUserID = respQueryOwnership.GetId()
	}
	return &ScheduledSearchV2{
		ID:                          respScheduledSearch.GetId(),
		Name:                        respScheduledSearch.GetName(),
		Description:                 respScheduledSearch.GetDescription(),
		QueryString:                 respScheduledSearch.GetQueryString(),
		SearchIntervalSeconds:       respScheduledSearch.GetSearchIntervalSeconds(),
		SearchIntervalOffsetSeconds: respScheduledSearch.GetSearchIntervalOffsetSeconds(),
		MaxWaitTimeSeconds:          respScheduledSearch.GetMaxWaitTimeSeconds(),
		TimeZone:                    respScheduledSearch.GetTimeZone(),
		Schedule:                    respScheduledSearch.GetSchedule(),
		BackfillLimitV2:             respScheduledSearch.GetBackfillLimitV2(),
		Enabled:                     respScheduledSearch.GetEnabled(),
		ActionNames:                 actionNames,
		OwnershipRunAsID:            runAsUserID,
		Labels:                      respScheduledSearch.GetLabels(),
		QueryTimestampType:          string(respScheduledSearch.QueryTimestampType),
		QueryOwnershipType:          string(queryOwnershipToQueryOwnershipType(respScheduledSearch.GetQueryOwnership())),
	}, nil
}

func (a *ScheduledSearchesV2) Delete(searchDomainName, scheduledSearchID string) error {
	if searchDomainName == "" {
		return fmt.Errorf("searchdomainName is empty")
	}
	if scheduledSearchID == "" {
		return fmt.Errorf("scheduledSearchID is empty")
	}

	_, err := humiographql.DeleteScheduledSearchV2ByID(context.Background(), a.client, searchDomainName, scheduledSearchID)
	return err
}
