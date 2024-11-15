package api

import (
	"context"
	"fmt"

	"github.com/humio/cli/internal/api/humiographql"
)

type AggregateAlert struct {
	ID                    string `yaml:"-"`
	Name                  string
	Description           *string
	QueryString           string   `yaml:"queryString"`
	SearchIntervalSeconds int64    `yaml:"searchIntervalSeconds"`
	ActionNames           []string `yaml:"actionNames"`
	Labels                []string
	Enabled               bool
	ThrottleField         *string `yaml:"throttleField"`
	ThrottleTimeSeconds   int64   `yaml:"throttleTimeSeconds"`
	QueryOwnershipType    string  `yaml:"queryOwnershipType"`
	TriggerMode           string  `yaml:"triggerMode"`
	QueryTimestampType    string  `yaml:"queryTimestampType"`
	OwnershipRunAsID      string  `yaml:"ownershipRunAsID"`
}

type AggregateAlerts struct {
	client *Client
}

func (c *Client) AggregateAlerts() *AggregateAlerts { return &AggregateAlerts{client: c} }

func (a *AggregateAlerts) List(searchDomainName string) ([]AggregateAlert, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("searchDomainName must not be empty")
	}

	resp, err := humiographql.ListAggregateAlerts(context.Background(), a.client, searchDomainName)
	if err != nil {
		return nil, err
	}
	respSearchDomain := resp.GetSearchDomain()
	respAggregateAlerts := respSearchDomain.GetAggregateAlerts()
	aggregateAlerts := make([]AggregateAlert, len(respAggregateAlerts))
	for idx, aggregateAlert := range respAggregateAlerts {
		actionNames := make([]string, len(aggregateAlert.GetActions()))
		for kdx, action := range aggregateAlert.GetActions() {
			actionNames[kdx] = action.GetName()
		}
		aggregateAlerts[idx] = AggregateAlert{
			ID:                    aggregateAlert.GetId(),
			Name:                  aggregateAlert.GetName(),
			Description:           aggregateAlert.GetDescription(),
			QueryString:           aggregateAlert.GetQueryString(),
			SearchIntervalSeconds: aggregateAlert.GetSearchIntervalSeconds(),
			ActionNames:           actionNames,
			Labels:                aggregateAlert.GetLabels(),
			Enabled:               aggregateAlert.GetEnabled(),
			ThrottleField:         aggregateAlert.ThrottleField,
			ThrottleTimeSeconds:   aggregateAlert.GetThrottleTimeSeconds(),
			QueryOwnershipType:    string(queryOwnershipToQueryOwnershipType(aggregateAlert.GetQueryOwnership())),
			TriggerMode:           string(aggregateAlert.GetTriggerMode()),
			QueryTimestampType:    string(aggregateAlert.GetQueryTimestampType()),
			OwnershipRunAsID:      aggregateAlert.GetQueryOwnership().GetId(),
		}
	}
	return aggregateAlerts, nil
}

func (a *AggregateAlerts) Update(searchDomainName string, updatedAggregateAlert *AggregateAlert) (*AggregateAlert, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("viewName must not be empty")
	}

	if updatedAggregateAlert == nil {
		return nil, fmt.Errorf("updatedAggregateAlert must not be nil")
	}

	if updatedAggregateAlert.ID == "" {
		return nil, fmt.Errorf("updatedAggregateAlert must have non-empty ID")
	}

	var ownershipRunAsID *string
	if humiographql.QueryOwnershipType(updatedAggregateAlert.QueryOwnershipType) == humiographql.QueryOwnershipTypeUser {
		ownershipRunAsID = &updatedAggregateAlert.OwnershipRunAsID
	}

	resp, err := humiographql.UpdateAggregateAlert(
		context.Background(),
		a.client,
		searchDomainName,
		updatedAggregateAlert.ID,
		updatedAggregateAlert.Name,
		updatedAggregateAlert.Description,
		updatedAggregateAlert.QueryString,
		updatedAggregateAlert.SearchIntervalSeconds,
		updatedAggregateAlert.ActionNames,
		updatedAggregateAlert.Labels,
		updatedAggregateAlert.Enabled,
		ownershipRunAsID,
		updatedAggregateAlert.ThrottleField,
		updatedAggregateAlert.ThrottleTimeSeconds,
		humiographql.TriggerMode(updatedAggregateAlert.TriggerMode),
		humiographql.QueryTimestampType(updatedAggregateAlert.QueryTimestampType),
		humiographql.QueryOwnershipType(updatedAggregateAlert.QueryOwnershipType),
	)
	if err != nil {
		return nil, err
	}

	respAggregateAlert := resp.GetUpdateAggregateAlert()
	actionNames := make([]string, len(respAggregateAlert.GetActions()))
	for kdx, action := range respAggregateAlert.GetActions() {
		actionNames[kdx] = action.GetName()
	}
	return &AggregateAlert{
		ID:                    respAggregateAlert.GetId(),
		Name:                  respAggregateAlert.GetName(),
		Description:           respAggregateAlert.GetDescription(),
		QueryString:           respAggregateAlert.GetQueryString(),
		SearchIntervalSeconds: respAggregateAlert.GetSearchIntervalSeconds(),
		ActionNames:           actionNames,
		Labels:                respAggregateAlert.GetLabels(),
		Enabled:               respAggregateAlert.GetEnabled(),
		ThrottleField:         respAggregateAlert.ThrottleField,
		ThrottleTimeSeconds:   respAggregateAlert.GetThrottleTimeSeconds(),
		QueryOwnershipType:    string(queryOwnershipToQueryOwnershipType(respAggregateAlert.GetQueryOwnership())),
		TriggerMode:           string(respAggregateAlert.GetTriggerMode()),
		QueryTimestampType:    string(respAggregateAlert.GetQueryTimestampType()),
		OwnershipRunAsID:      respAggregateAlert.GetQueryOwnership().GetId(),
	}, nil
}

func (a *AggregateAlerts) Create(searchDomainName string, newAggregateAlert *AggregateAlert) (*AggregateAlert, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("viewName must not be empty")
	}

	if newAggregateAlert == nil {
		return nil, fmt.Errorf("newAggregateAlert must not be nil")
	}

	var ownershipRunAsID *string
	if humiographql.QueryOwnershipType(newAggregateAlert.QueryOwnershipType) == humiographql.QueryOwnershipTypeUser {
		ownershipRunAsID = &newAggregateAlert.OwnershipRunAsID
	}

	resp, err := humiographql.CreateAggregateAlert(
		context.Background(),
		a.client,
		searchDomainName,
		newAggregateAlert.Name,
		newAggregateAlert.Description,
		newAggregateAlert.QueryString,
		newAggregateAlert.SearchIntervalSeconds,
		newAggregateAlert.ActionNames,
		newAggregateAlert.Labels,
		newAggregateAlert.Enabled,
		ownershipRunAsID,
		newAggregateAlert.ThrottleField,
		newAggregateAlert.ThrottleTimeSeconds,
		humiographql.TriggerMode(newAggregateAlert.TriggerMode),
		humiographql.QueryTimestampType(newAggregateAlert.QueryTimestampType),
		humiographql.QueryOwnershipType(newAggregateAlert.QueryOwnershipType),
	)
	if err != nil {
		return nil, err
	}

	respAggregateAlert := resp.GetCreateAggregateAlert()
	actionNames := make([]string, len(respAggregateAlert.GetActions()))
	for kdx, action := range respAggregateAlert.GetActions() {
		actionNames[kdx] = action.GetName()
	}
	return &AggregateAlert{
		ID:                    respAggregateAlert.GetId(),
		Name:                  respAggregateAlert.GetName(),
		Description:           respAggregateAlert.GetDescription(),
		QueryString:           respAggregateAlert.GetQueryString(),
		SearchIntervalSeconds: respAggregateAlert.GetSearchIntervalSeconds(),
		ActionNames:           actionNames,
		Labels:                respAggregateAlert.GetLabels(),
		Enabled:               respAggregateAlert.GetEnabled(),
		ThrottleField:         respAggregateAlert.ThrottleField,
		ThrottleTimeSeconds:   respAggregateAlert.GetThrottleTimeSeconds(),
		QueryOwnershipType:    string(queryOwnershipToQueryOwnershipType(respAggregateAlert.GetQueryOwnership())),
		TriggerMode:           string(respAggregateAlert.GetTriggerMode()),
		QueryTimestampType:    string(respAggregateAlert.GetQueryTimestampType()),
		OwnershipRunAsID:      respAggregateAlert.GetQueryOwnership().GetId(),
	}, nil
}

func (a *AggregateAlerts) Delete(searchDomainName, aggregateAlertID string) error {
	if searchDomainName == "" {
		return fmt.Errorf("viewName must not be empty")
	}

	if aggregateAlertID == "" {
		return fmt.Errorf("aggregateAlertID is empty")
	}

	_, err := humiographql.DeleteAggregateAlert(context.Background(), a.client, searchDomainName, aggregateAlertID)
	return err
}

func (a *AggregateAlerts) Get(searchDomainName string, aggregateAlertID string) (*AggregateAlert, error) {
	if searchDomainName == "" {
		return nil, fmt.Errorf("viewName must not be empty")
	}

	if aggregateAlertID == "" {
		return nil, fmt.Errorf("aggregateAlertID must not be empty")
	}

	resp, err := humiographql.GetAggregateAlertByID(context.Background(), a.client, searchDomainName, aggregateAlertID)
	if err != nil {
		return nil, err
	}
	respSearchDomain := resp.GetSearchDomain()
	respAggregateAlert := respSearchDomain.GetAggregateAlert()
	actionNames := make([]string, len(respAggregateAlert.GetActions()))
	for kdx, action := range respAggregateAlert.GetActions() {
		actionNames[kdx] = action.GetName()
	}
	return &AggregateAlert{
		ID:                    respAggregateAlert.GetId(),
		Name:                  respAggregateAlert.GetName(),
		Description:           respAggregateAlert.GetDescription(),
		QueryString:           respAggregateAlert.GetQueryString(),
		SearchIntervalSeconds: respAggregateAlert.GetSearchIntervalSeconds(),
		ActionNames:           actionNames,
		Labels:                respAggregateAlert.GetLabels(),
		Enabled:               respAggregateAlert.GetEnabled(),
		ThrottleField:         respAggregateAlert.ThrottleField,
		ThrottleTimeSeconds:   respAggregateAlert.GetThrottleTimeSeconds(),
		QueryOwnershipType:    string(queryOwnershipToQueryOwnershipType(respAggregateAlert.GetQueryOwnership())),
		TriggerMode:           string(respAggregateAlert.GetTriggerMode()),
		QueryTimestampType:    string(respAggregateAlert.GetQueryTimestampType()),
		OwnershipRunAsID:      respAggregateAlert.GetQueryOwnership().GetId(),
	}, nil
}
