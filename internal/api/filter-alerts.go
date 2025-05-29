package api

import (
	"context"
	"fmt"

	"github.com/humio/cli/internal/api/humiographql"
)

type FilterAlert struct {
	ID                  string `yaml:"-"`
	Name                string
	Description         *string
	QueryString         string   `yaml:"queryString"`
	ActionNames         []string `yaml:"actionNames"`
	Labels              []string
	Enabled             bool
	QueryOwnershipType  string  `yaml:"queryOwnershipType"`
	ThrottleTimeSeconds *int64  `yaml:"throttleTimeSeconds"`
	ThrottleField       *string `yaml:"throttleField"`
	OwnershipRunAsID    string  `yaml:"ownershipRunAsID"`
}

type FilterAlerts struct {
	client *Client
}

func (c *Client) FilterAlerts() *FilterAlerts { return &FilterAlerts{client: c} }

func (fa *FilterAlerts) List(searchDomainName string) ([]FilterAlert, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("searchDomainName must not be empty")
	}

	resp, err := humiographql.ListFilterAlerts(context.Background(), fa.client, searchDomainName)
	if err != nil {
		return nil, err
	}
	respSearchDomain := resp.GetSearchDomain()
	respFilterAlerts := respSearchDomain.GetFilterAlerts()
	filterAlerts := make([]FilterAlert, len(respFilterAlerts))
	for idx, filterAlert := range respFilterAlerts {
		actionNames := make([]string, len(filterAlert.GetActions()))
		for kdx, action := range filterAlert.GetActions() {
			actionNames[kdx] = action.GetName()
		}
		filterAlerts[idx] = FilterAlert{
			ID:                  filterAlert.GetId(),
			Name:                filterAlert.GetName(),
			Description:         filterAlert.GetDescription(),
			QueryString:         filterAlert.GetQueryString(),
			ActionNames:         actionNames,
			Labels:              filterAlert.GetLabels(),
			Enabled:             filterAlert.GetEnabled(),
			ThrottleField:       filterAlert.GetThrottleField(),
			ThrottleTimeSeconds: filterAlert.GetThrottleTimeSeconds(),
			QueryOwnershipType:  string(queryOwnershipToQueryOwnershipType(filterAlert.GetQueryOwnership())),
			OwnershipRunAsID:    filterAlert.GetQueryOwnership().GetId(),
		}
	}
	return filterAlerts, nil
}

func (fa *FilterAlerts) Create(searchDomainName string, newFilterAlert *FilterAlert) (*FilterAlert, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("searchDomainName must not be empty")
	}

	if newFilterAlert == nil {
		return nil, fmt.Errorf("newFilterAlert must not be nil")
	}

	var ownershipRunAsID *string
	if humiographql.QueryOwnershipType(newFilterAlert.QueryOwnershipType) == humiographql.QueryOwnershipTypeUser {
		ownershipRunAsID = &newFilterAlert.OwnershipRunAsID
	}

	resp, err := humiographql.CreateFilterAlert(
		context.Background(),
		fa.client,
		searchDomainName,
		newFilterAlert.Name,
		newFilterAlert.Description,
		newFilterAlert.QueryString,
		newFilterAlert.ActionNames,
		newFilterAlert.Labels,
		newFilterAlert.Enabled,
		ownershipRunAsID,
		newFilterAlert.ThrottleField,
		newFilterAlert.ThrottleTimeSeconds,
		humiographql.QueryOwnershipType(newFilterAlert.QueryOwnershipType),
	)
	if err != nil {
		return nil, err
	}

	respFilterAlert := resp.GetCreateFilterAlert()
	actionNames := make([]string, len(respFilterAlert.GetActions()))
	for kdx, action := range respFilterAlert.GetActions() {
		actionNames[kdx] = action.GetName()
	}
	return &FilterAlert{
		ID:                  respFilterAlert.GetId(),
		Name:                respFilterAlert.GetName(),
		Description:         respFilterAlert.GetDescription(),
		QueryString:         respFilterAlert.GetQueryString(),
		ActionNames:         actionNames,
		Labels:              respFilterAlert.GetLabels(),
		Enabled:             respFilterAlert.GetEnabled(),
		ThrottleField:       respFilterAlert.ThrottleField,
		ThrottleTimeSeconds: respFilterAlert.GetThrottleTimeSeconds(),
		QueryOwnershipType:  string(queryOwnershipToQueryOwnershipType(respFilterAlert.GetQueryOwnership())),
		OwnershipRunAsID:    respFilterAlert.GetQueryOwnership().GetId(),
	}, nil
}

func (fa *FilterAlerts) Delete(searchDomainName, filterAlertID string) error {
	if filterAlertID == "" {
		return fmt.Errorf("filterAlertID is empty")
	}

	_, err := humiographql.DeleteFilterAlert(context.Background(), fa.client, searchDomainName, filterAlertID)
	return err
}
