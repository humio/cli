package api

import (
	"context"
	"fmt"

	"github.com/humio/cli/internal/api/humiographql"
)

type Alert struct {
	ID                 string `yaml:"-"`
	Name               string
	QueryString        string  `yaml:"queryString"`
	QueryStart         string  `yaml:"queryStart"`
	ThrottleField      *string `yaml:"throttleField"`
	TimeOfLastTrigger  *int64  `yaml:"timeOfLastTrigger"`
	IsStarred          bool    `yaml:"isStarred"`
	Description        *string
	ThrottleTimeMillis int64 `yaml:"throttleTimeMillis"`
	Enabled            bool
	Actions            []string
	Labels             []string
	LastError          *string `yaml:"lastError"`
	RunAsUserID        string  `yaml:"runAsUserID"`
	QueryOwnershipType string  `yaml:"queryOwnershipType"`
}

type Alerts struct {
	client *Client
}

func (c *Client) Alerts() *Alerts { return &Alerts{client: c} }

func (a *Alerts) List(searchDomainName string) ([]Alert, error) {
	resp, err := humiographql.ListAlerts(context.Background(), a.client, searchDomainName)
	if err != nil {
		return nil, err
	}
	respSearchDomain := resp.GetSearchDomain()
	respAlerts := respSearchDomain.GetAlerts()
	alerts := make([]Alert, len(respAlerts))
	for idx, alert := range respAlerts {
		respOwnership := alert.GetQueryOwnership()
		runAsUserID := ""
		if respOwnership != nil {
			runAsUserID = respOwnership.GetId()
		}
		alerts[idx] = Alert{
			ID:                 alert.GetId(),
			Name:               alert.GetName(),
			QueryString:        alert.GetQueryString(),
			QueryStart:         alert.GetQueryStart(),
			ThrottleField:      alert.GetThrottleField(),
			TimeOfLastTrigger:  alert.GetTimeOfLastTrigger(),
			IsStarred:          alert.GetIsStarred(),
			Description:        alert.GetDescription(),
			ThrottleTimeMillis: alert.GetThrottleTimeMillis(),
			Enabled:            alert.GetEnabled(),
			Actions:            alert.GetActions(),
			Labels:             alert.GetLabels(),
			LastError:          alert.GetLastError(),
			RunAsUserID:        runAsUserID,
			QueryOwnershipType: string(queryOwnershipToQueryOwnershipType(alert.GetQueryOwnership())),
		}
	}
	return alerts, nil
}

func (a *Alerts) Update(searchDomainName string, newAlert *Alert) (*Alert, error) {
	if newAlert == nil {
		return nil, fmt.Errorf("newAlert must not be nil")
	}

	if newAlert.ID == "" {
		return nil, fmt.Errorf("newAlert must have non-empty newAlert id")
	}

	queryOwnershipType := humiographql.QueryOwnershipType(newAlert.QueryOwnershipType)

	var ownershipRunAsID *string
	if queryOwnershipType == humiographql.QueryOwnershipTypeUser {
		ownershipRunAsID = &newAlert.RunAsUserID
	}

	resp, err := humiographql.UpdateAlert(
		context.Background(),
		a.client,
		searchDomainName,
		newAlert.ID,
		newAlert.Name,
		newAlert.Description,
		newAlert.QueryString,
		newAlert.QueryStart,
		newAlert.ThrottleTimeMillis,
		newAlert.Enabled,
		newAlert.Actions,
		newAlert.Labels,
		ownershipRunAsID,
		&queryOwnershipType,
		newAlert.ThrottleField,
	)
	if err != nil {
		return nil, err
	}

	respUpdate := resp.GetUpdateAlert()
	respQueryOwnership := respUpdate.GetQueryOwnership()
	runAsUserID := ""
	if respQueryOwnership != nil {
		runAsUserID = respQueryOwnership.GetId()
	}
	return &Alert{
		ID:                 respUpdate.GetId(),
		Name:               respUpdate.GetName(),
		QueryString:        respUpdate.GetQueryString(),
		QueryStart:         respUpdate.GetQueryStart(),
		ThrottleField:      respUpdate.GetThrottleField(),
		TimeOfLastTrigger:  respUpdate.GetTimeOfLastTrigger(),
		IsStarred:          respUpdate.GetIsStarred(),
		Description:        respUpdate.GetDescription(),
		ThrottleTimeMillis: respUpdate.GetThrottleTimeMillis(),
		Enabled:            respUpdate.GetEnabled(),
		Actions:            respUpdate.GetActions(),
		Labels:             respUpdate.GetLabels(),
		LastError:          respUpdate.LastError,
		RunAsUserID:        runAsUserID,
		QueryOwnershipType: *respQueryOwnership.GetTypename(),
	}, nil
}

func (a *Alerts) Add(searchDomainName string, newAlert *Alert) (*Alert, error) {
	if newAlert == nil {
		return nil, fmt.Errorf("newAlert must not be nil")
	}

	queryOwnershipType := humiographql.QueryOwnershipType(newAlert.QueryOwnershipType)

	var ownershipRunAsID *string
	if queryOwnershipType == humiographql.QueryOwnershipTypeUser {
		ownershipRunAsID = &newAlert.RunAsUserID
	}

	resp, err := humiographql.CreateAlert(
		context.Background(),
		a.client,
		searchDomainName,
		newAlert.Name,
		newAlert.Description,
		newAlert.QueryString,
		newAlert.QueryStart,
		newAlert.ThrottleTimeMillis,
		&newAlert.Enabled,
		newAlert.Actions,
		newAlert.Labels,
		ownershipRunAsID,
		&queryOwnershipType,
		newAlert.ThrottleField,
	)
	if err != nil {
		return nil, err
	}

	respUpdate := resp.GetCreateAlert()
	respQueryOwnership := respUpdate.GetQueryOwnership()
	runAsUserID := ""
	if respQueryOwnership != nil {
		runAsUserID = respQueryOwnership.GetId()
	}
	return &Alert{
		ID:                 respUpdate.GetId(),
		Name:               respUpdate.GetName(),
		QueryString:        respUpdate.GetQueryString(),
		QueryStart:         respUpdate.GetQueryStart(),
		ThrottleField:      respUpdate.GetThrottleField(),
		TimeOfLastTrigger:  respUpdate.GetTimeOfLastTrigger(),
		IsStarred:          respUpdate.GetIsStarred(),
		Description:        respUpdate.GetDescription(),
		ThrottleTimeMillis: respUpdate.GetThrottleTimeMillis(),
		Enabled:            respUpdate.GetEnabled(),
		Actions:            respUpdate.GetActions(),
		Labels:             respUpdate.GetLabels(),
		LastError:          respUpdate.LastError,
		RunAsUserID:        runAsUserID,
		QueryOwnershipType: *respQueryOwnership.GetTypename(),
	}, nil
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

func (a *Alerts) Delete(searchDomainName, alertName string) error {
	actions, err := a.List(searchDomainName)
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
		return AlertNotFound(alertName)
	}

	_, err = humiographql.DeleteAlert(context.Background(), a.client, searchDomainName, alertName)
	return err
}
