package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const NotifierTypeEmail = "EmailNotifier"
const NotifierTypeHumioRepo = "HumioRepoNotifier"
const NotifierTypeOpsGenie = "OpsGenieNotifier"
const NotifierTypePagerDuty = "PagerDutyNotifier"
const NotifierTypeSlack = "SlackNotifier"
const NotifierTypeSlackPostMessage = "SlackPostMessageNotifier"
const NotifierTypeVictorOps = "VictorOpsNotifier"
const NotifierTypeWebHook = "WebHookNotifier"

type Notifiers struct {
	client *Client
}

type Notifier struct {
	Entity     string                 `json:"entity"`
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
}

func (c *Client) Notifiers() *Notifiers { return &Notifiers{client: c} }

func (n *Notifiers) List(viewName string) ([]Notifier, error) {
	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers", viewName)

	res, err := n.client.HTTPRequest(http.MethodGet, url, nil)
	if err != nil {
		return []Notifier{}, err
	}

	return n.unmarshalToNotifierList(res)
}

func (n *Notifiers) Update(viewName string, notifier *Notifier) (*Notifier, error) {
	if notifier.ID == "" {
		existingID, err := n.convertNotifierNameToID(viewName, notifier.Name)
		if err != nil {
			return nil, fmt.Errorf("could not convert notifier name to id: %w", err)
		}
		notifier.ID = existingID
	}

	jsonStr, err := n.marshalToJSON(notifier)
	if err != nil {
		return nil, fmt.Errorf("unable to convert notifier to json string: %w", err)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers/%s", viewName, notifier.ID)

	res, err := n.client.HTTPRequest(http.MethodPut, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	return n.unmarshalToNotifier(res)
}

func (n *Notifiers) Add(viewName string, notifier *Notifier, force bool) (*Notifier, error) {
	nameAlreadyInUse, err := n.notifierNameInUse(viewName, notifier.Name)
	if err != nil {
		return nil, fmt.Errorf("could not determine if notifier name is in use: %w", err)
	}
	if nameAlreadyInUse {
		if !force {
			return nil, fmt.Errorf("notifier with name %s already exists", notifier.Name)
		}
		return n.Update(viewName, notifier)
	}

	jsonStr, err := n.marshalToJSON(notifier)
	if err != nil {
		return nil, fmt.Errorf("unable to convert notifier to json string: %w", err)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers/", viewName)

	res, err := n.client.HTTPRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	return n.unmarshalToNotifier(res)
}

func (n *Notifiers) Get(viewName, notifierName string) (*Notifier, error) {
	notifierID, err := n.convertNotifierNameToID(viewName, notifierName)
	if err != nil {
		return nil, fmt.Errorf("could not find a notifier in view %s with name: %s", viewName, notifierName)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers/%s", viewName, notifierID)

	res, err := n.client.HTTPRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return n.unmarshalToNotifier(res)
}

func (n *Notifiers) GetByID(viewName, notifierID string) (*Notifier, error) {
	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers/%s", viewName, notifierID)

	res, err := n.client.HTTPRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return n.unmarshalToNotifier(res)
}

func (n *Notifiers) Delete(viewName, notifierName string) error {
	notifierID, err := n.convertNotifierNameToID(viewName, notifierName)
	if err != nil {
		return fmt.Errorf("could not find a notifier in view %s with name: %s", viewName, notifierName)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers/%s", viewName, notifierID)

	res, err := n.client.HTTPRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	if err != nil || res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete notifier in view %s with id %s, got: %w", viewName, notifierID, err)
	}

	return nil
}

func (n *Notifiers) marshalToJSON(notifier *Notifier) ([]byte, error) {
	jsonStr, err := json.Marshal(notifier)
	if err != nil {
		return nil, fmt.Errorf("unable to convert notifier to json string: %w", err)
	}
	return jsonStr, nil
}

func (n *Notifiers) unmarshalToNotifierList(res *http.Response) ([]Notifier, error) {
	notifiers := []Notifier{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return notifiers, err
	}

	if err = json.Unmarshal(body, &notifiers); err != nil {
		return notifiers, err
	}
	return notifiers, nil
}

func (n *Notifiers) unmarshalToNotifier(res *http.Response) (*Notifier, error) {
	notifier := Notifier{}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &notifier, err
	}

	if err = json.Unmarshal(body, &notifier); err != nil {
		return &notifier, fmt.Errorf("error in json response: %w. response: %v", err, string(body))
	}

	return &notifier, nil
}

func (n *Notifiers) convertNotifierNameToID(viewName, notifierName string) (string, error) {
	listOfNotifiers, err := n.List(viewName)
	if err != nil {
		return "", fmt.Errorf("could not list all notifiers for view %s: %w", viewName, err)
	}
	for _, v := range listOfNotifiers {
		if v.Name == notifierName {
			return v.ID, nil
		}
	}
	return "", fmt.Errorf("could not find a notifier in view %s with name: %s", viewName, notifierName)
}

func (n *Notifiers) notifierNameInUse(viewName, notifierName string) (bool, error) {
	listOfNotifiers, err := n.List(viewName)
	if err != nil {
		return true, fmt.Errorf("could not list all notifiers for view %s: %w", viewName, err)
	}
	for _, v := range listOfNotifiers {
		if v.Name == notifierName {
			return true, nil
		}
	}
	return false, nil
}
