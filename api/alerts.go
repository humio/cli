package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type HumioQuery struct {
	QueryString string `yaml:"queryString" json:"queryString"`
	Start       string `yaml:"start"       json:"start"`
	End         string `yaml:"end"         json:"end"`
	IsLive      bool   `yaml:"isLive"      json:"isLive"`
}

type Alert struct {
	ID                 string     `yaml:"-"                     json:"id"`
	Name               string     `yaml:"name"                  json:"name"`
	Query              HumioQuery `yaml:"query"                 json:"query"`
	Description        string     `yaml:"description,omitempty" json:"description"`
	ThrottleTimeMillis int        `yaml:"throttleTimeMillis"    json:"throttleTimeMillis"`
	Silenced           bool       `yaml:"silenced"              json:"silenced"`
	Notifiers          []string   `yaml:"notifiers"             json:"notifiers"`
	LinkURL            string     `yaml:"linkURL"               json:"linkURL"`
	Labels             []string   `yaml:"labels,omitempty"      json:"labels,omitempty"`
}

type Alerts struct {
	client *Client
}

func (c *Client) Alerts() *Alerts { return &Alerts{client: c} }

func (a *Alerts) List(viewName string) ([]Alert, error) {
	url := fmt.Sprintf("api/v1/repositories/%s/alerts", viewName)

	res, err := a.client.HTTPRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return a.unmarshalToAlertList(res)
}

func (a *Alerts) Update(viewName string, alert *Alert) (*Alert, error) {
	existingID, err := a.convertAlertNameToID(viewName, alert.Name)
	if err != nil {
		return nil, fmt.Errorf("could not convert alert name to id: %v", err)
	}

	jsonStr, err := a.marshalToJSON(alert)
	if err != nil {
		return nil, fmt.Errorf("unable to convert alert to json string: %v", err)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alerts/%s", viewName, existingID)

	res, postErr := a.client.HTTPRequest(http.MethodPut, url, bytes.NewBuffer(jsonStr))
	if postErr != nil {
		return nil, fmt.Errorf("could not add alert in view %s with name %s, got: %v", viewName, alert.Name, postErr)
	}
	return a.unmarshalToAlert(res)
}

func (a *Alerts) Add(viewName string, alert *Alert, updateExisting bool) (*Alert, error) {
	nameAlreadyInUse, err := a.alertNameInUse(viewName, alert.Name)
	if err != nil {
		return nil, fmt.Errorf("could not determine if alert name is in use: %v", err)
	}

	if nameAlreadyInUse {
		if updateExisting == false {
			return nil, fmt.Errorf("alert with name %s already exists", alert.Name)
		}
		return a.Update(viewName, alert)
	}

	jsonStr, err := a.marshalToJSON(alert)
	if err != nil {
		return nil, fmt.Errorf("unable to convert alert to json string: %v", err)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alerts/", viewName)

	res, err := a.client.HTTPRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("could not add alert in view %s with name %s, got: %v", viewName, alert.Name, err)
	}

	return a.unmarshalToAlert(res)
}

func (a *Alerts) Get(viewName, alertName string) (*Alert, error) {
	alertID, err := a.convertAlertNameToID(viewName, alertName)
	if err != nil {
		return nil, fmt.Errorf("could not find a notifier in view %s with name: %s", viewName, alertName)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alerts/%s", viewName, alertID)

	res, err := a.client.HTTPRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get alert with id %s, got: %v", alertID, err)
	}

	return a.unmarshalToAlert(res)
}

func (a *Alerts) Delete(viewName, alertName string) error {
	alertID, err := a.convertAlertNameToID(viewName, alertName)
	if err != nil {
		return fmt.Errorf("could not find a notifier in view %s with name: %s", viewName, alertName)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alerts/%s", viewName, alertID)

	res, err := a.client.HTTPRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil || res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete alert in view %s with id %s, got: %v", viewName, alertID, err)
	}
	return nil
}

func (a *Alerts) marshalToJSON(alert *Alert) ([]byte, error) {
	// Humio requires notifiers to be specified even if no notifier is desired
	if alert.Notifiers == nil {
		alert.Notifiers = []string{}
	}

	jsonStr, err := json.Marshal(alert)
	if err != nil {
		return nil, fmt.Errorf("unable to convert alert to json string: %v", err)
	}
	return jsonStr, nil
}

func (a *Alerts) unmarshalToAlertList(res *http.Response) ([]Alert, error) {
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	alertList := []Alert{}
	jsonErr := json.Unmarshal(body, &alertList)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return alertList, nil
}

func (a *Alerts) unmarshalToAlert(res *http.Response) (*Alert, error) {
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	alert := Alert{}
	jsonErr := json.Unmarshal(body, &alert)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return &alert, nil
}

func (a *Alerts) convertAlertNameToID(viewName, alertName string) (string, error) {
	listOfAlerts, err := a.List(viewName)
	if err != nil {
		return "", fmt.Errorf("could not list all alerts for view %s: %v", viewName, err)
	}
	for _, v := range listOfAlerts {
		if v.Name == alertName {
			return v.ID, nil
		}
	}
	return "", fmt.Errorf("could not find an alert in view %s with name: %s", viewName, alertName)
}

func (a *Alerts) alertNameInUse(viewName, alertName string) (bool, error) {
	listOfAlerts, err := a.List(viewName)
	if err != nil {
		return true, fmt.Errorf("could not list all alerts for view %s: %v", viewName, err)
	}
	for _, v := range listOfAlerts {
		if v.Name == alertName {
			return true, nil
		}
	}
	return false, nil
}
