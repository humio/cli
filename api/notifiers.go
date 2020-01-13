package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const NotifierTypeEmail = "EmailNotifier"
const NotifierTypeOpsGenie = "OpsGenieNotifier"
const NotifierTypePagerDuty = "PagerDutyNotifier"
const NotifierTypeSlack = "SlackNotifier"
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

func (n *Notifiers) List(view string) ([]Notifier, error) {
	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers", view)

	res, err := n.client.HttpGET(url)
	if err != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	notifiers := []Notifier{}
	jsonErr := json.Unmarshal(body, &notifiers)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return notifiers, nil
}

func (n *Notifiers) Update(viewName string, notifier *Notifier) (*Notifier, error) {
	existingID, err := n.convertNotifierNameToID(viewName, notifier.Name)
	if err != nil {
		return nil, fmt.Errorf("could not convert notifier name to id: %v", err)
	}
	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers/%s", viewName, existingID)

	jsonStr, err := json.Marshal(notifier)
	if err != nil {
		return nil, fmt.Errorf("unable to convert notifier to json string: %v", err)
	}

	res, err := n.client.HttpPUT(url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	fmt.Println(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	resNotifier := Notifier{}
	jsonErr := json.Unmarshal(body, &resNotifier)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return &resNotifier, nil
}

func (n *Notifiers) Add(viewName string, notifier *Notifier, force bool) (*Notifier, error) {
	url := fmt.Sprintf("%sapi/v1/repositories/%s/alertnotifiers/", n.client.Address(), viewName)

	nameAlreadyInUse, err := n.notifierNameInUse(viewName, notifier.Name)
	if err != nil {
		return nil, fmt.Errorf("could not determine if notifier name is in use: %v", err)
	}
	if nameAlreadyInUse {
		if force == false {
			return nil, fmt.Errorf("notifier with name %s already exists", notifier.Name)
		}
		return n.Update(viewName, notifier)
	}

	jsonStr, err := json.Marshal(notifier)
	if err != nil {
		return nil, fmt.Errorf("unable to convert notifier to json string: %v", err)
	}

	res, err := n.client.HttpPOST(url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	fmt.Println(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	resNotifier := Notifier{}
	jsonErr := json.Unmarshal(body, &resNotifier)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return &resNotifier, nil
}

func (n *Notifiers) Get(view, name string) (*Notifier, error) {
	notifierID, err := n.convertNotifierNameToID(view, name)
	if err != nil {
		return nil, fmt.Errorf("could not find a notifier in view %s with name: %s", view, name)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers/%s", view, notifierID)

	res, err := n.client.HttpGET(url)
	if err != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	resNotifier := Notifier{}
	jsonErr := json.Unmarshal(body, &resNotifier)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return &resNotifier, nil
}

func (n *Notifiers) GetByID(view, notifierID string) (*Notifier, error) {
	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers/%s", view, notifierID)

	res, err := n.client.HttpGET(url)
	if err != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	resNotifier := Notifier{}
	jsonErr := json.Unmarshal(body, &resNotifier)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return &resNotifier, nil
}

func (n *Notifiers) Delete(viewName, name string) error {
	notifierID, err := n.convertNotifierNameToID(viewName, name)
	if err != nil {
		return fmt.Errorf("could not find a notifier in view %s with name: %s", viewName, name)
	}

	url := fmt.Sprintf("api/v1/repositories/%s/alertnotifiers/%s", viewName, notifierID)

	res, err := n.client.HttpDELETE(url)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil || res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("could not delete notifier in view %s with id %s, got: %v", viewName, notifierID, err)
	}
	return nil
}

func (n *Notifiers) convertNotifierNameToID(viewName, notifierName string) (string, error) {
	listOfNotifiers, err := n.List(viewName)
	if err != nil {
		return "", fmt.Errorf("could not list all notifiers for view %s: %v", viewName, err)
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
		return true, fmt.Errorf("could not list all notifiers for view %s: %v", viewName, err)
	}
	for _, v := range listOfNotifiers {
		if v.Name == notifierName {
			return true, nil
		}
	}
	return false, nil
}
