package api

import (
	"fmt"
	graphql "github.com/cli/shurcooL-graphql"
	"reflect"
)

const (
	ActionTypeEmail            string = "EmailAction"
	ActionTypeHumioRepo        string = "HumioRepoAction"
	ActionTypeOpsGenie         string = "OpsGenieAction"
	ActionTypePagerDuty        string = "PagerDutyAction"
	ActionTypeSlack            string = "SlackAction"
	ActionTypeSlackPostMessage string = "SlackPostMessageAction"
	ActionTypeVictorOps        string = "VictorOpsAction"
	ActionTypeWebhook          string = "WebhookAction"
)

type Actions struct {
	client *Client
}

type EmailAction struct {
	Recipients      []string `graphql:"emailRecipients: recipients" yaml:"recipients,omitempty" json:"recipients,omitempty"`
	SubjectTemplate string   `graphql:"emailSubjectTemplate: subjectTemplate" yaml:"subjectTemplate,omitempty" json:"subjectTemplate,omitempty"`
	BodyTemplate    string   `graphql:"emailBodyTemplate: bodyTemplate" yaml:"bodyTemplate,omitempty" json:"bodyTemplate,omitempty"`
	UseProxy        bool     `graphql:"emailUseProxy: useProxy" yaml:"useProxy,omitempty" json:"useProxy,omitempty"`
}

type HumioRepoAction struct {
	IngestToken string `graphql:"humioRepoIngestToken: ingestToken" yaml:"ingestToken,omitempty" json:"ingestToken,omitempty"`
}

type OpsGenieAction struct {
	ApiUrl   string `graphql:"opsGenieApiUrl: apiUrl" yaml:"apiUrl,omitempty" json:"apiUrl,omitempty"`
	GenieKey string `graphql:"opsGenieGenieKey: genieKey" yaml:"genieKey,omitempty" json:"genieKey,omitempty"`
	UseProxy bool   `graphql:"opsGenieUseProxy: useProxy" yaml:"useProxy,omitempty" json:"useProxy,omitempty"`
}

type PagerDutyAction struct {
	Severity   string `graphql:"pagerDutySeverity: severity" yaml:"severity,omitempty" json:"severity,omitempty"`
	RoutingKey string `graphql:"pagerDutyRoutingKey: routingKey" yaml:"routingKey,omitempty" json:"routingKey,omitempty"`
	UseProxy   bool   `graphql:"pagerDutyUseProxy: useProxy" yaml:"useProxy,omitempty" json:"useProxy,omitempty"`
}

type SlackFieldEntryInput struct {
	FieldName string `graphql:"fieldName" yaml:"fieldName" json:"fieldName"`
	Value     string `graphql:"value"     yaml:"value" json:"value"`
}

type SlackAction struct {
	Url      string                 `graphql:"slackUrl: url" yaml:"url,omitempty" json:"url,omitempty"`
	Fields   []SlackFieldEntryInput `graphql:"slackFields: fields" yaml:"fields,omitempty" json:"fields,omitempty"`
	UseProxy bool                   `graphql:"slackUseProxy: useProxy" yaml:"useProxy,omitempty" json:"useProxy,omitempty"`
}

type SlackPostMessageAction struct {
	ApiToken string                 `graphql:"slackPostMessageApiToken: apiToken" yaml:"apiToken,omitempty" json:"apiToken,omitempty"`
	Channels []string               `graphql:"slackPostMessageChannels: channels" yaml:"channels,omitempty" json:"channels,omitempty"`
	Fields   []SlackFieldEntryInput `graphql:"slackPostMessageFields: fields" yaml:"fields,omitempty" json:"fields,omitempty"`
	UseProxy bool                   `graphql:"slackPostMessageUseProxy: useProxy" yaml:"useProxy,omitempty" json:"useProxy,omitempty"`
}

type VictorOpsAction struct {
	MessageType string `graphql:"victorOpsMessageType: messageType" yaml:"messageType,omitempty" json:"messageType,omitempty"`
	NotifyUrl   string `graphql:"victorOpsNotifyUrl: notifyUrl" yaml:"notifyUrl,omitempty" json:"notifyUrl,omitempty"`
	UseProxy    bool   `graphql:"victorOpsUseProxy: useProxy" yaml:"useProxy,omitempty" json:"useProxy,omitempty"`
}

type HttpHeaderEntryInput struct {
	Header string `graphql:"header"  yaml:"header" json:"header"`
	Value  string `graphql:"value"   yaml:"value" json:"value"`
}

type WebhookAction struct {
	Method       string                 `graphql:"webhookMethod: method" yaml:"method,omitempty" json:"method,omitempty"`
	Url          string                 `graphql:"webhookUrl: url" yaml:"url,omitempty" json:"url,omitempty"`
	Headers      []HttpHeaderEntryInput `graphql:"webhookHeaders: headers" yaml:"headers,omitempty" json:"headers,omitempty"`
	BodyTemplate string                 `graphql:"webhookBodyTemplate: bodyTemplate" yaml:"bodyTemplate,omitempty" json:"bodyTemplate,omitempty"`
	IgnoreSSL    bool                   `graphql:"webhookIgnoreSSL: ignoreSSL" yaml:"ignoreSSL,omitempty" json:"ignoreSSL,omitempty"`
	UseProxy     bool                   `graphql:"webhookUseProxy: useProxy" yaml:"useProxy,omitempty" json:"useProxy,omitempty"`
}

type Action struct {
	Type string `graphql:"__typename" yaml:"type" json:"type"`
	ID   string `graphql:"id"         yaml:"-"    json:"id"`
	Name string `graphql:"name"       yaml:"name" json:"name"`

	EmailAction            EmailAction            `graphql:"... on EmailAction"            yaml:"emailAction,omitempty" json:"emailAction,omitempty"`
	HumioRepoAction        HumioRepoAction        `graphql:"... on HumioRepoAction"        yaml:"humioRepoAction,omitempty" json:"humioRepoAction,omitempty"`
	OpsGenieAction         OpsGenieAction         `graphql:"... on OpsGenieAction"         yaml:"opsGenieAction,omitempty" json:"opsGenieAction,omitempty"`
	PagerDutyAction        PagerDutyAction        `graphql:"... on PagerDutyAction"        yaml:"pagerDutyAction,omitempty" json:"pagerDutyAction,omitempty"`
	SlackAction            SlackAction            `graphql:"... on SlackAction"            yaml:"slackAction,omitempty" json:"slackAction,omitempty"`
	SlackPostMessageAction SlackPostMessageAction `graphql:"... on SlackPostMessageAction" yaml:"slackPostMessageAction,omitempty" json:"slackPostMessageAction,omitempty"`
	VictorOpsAction        VictorOpsAction        `graphql:"... on VictorOpsAction"        yaml:"victorOpsAction,omitempty" json:"victorOpsAction,omitempty"`
	WebhookAction          WebhookAction          `graphql:"... on WebhookAction"          yaml:"webhookAction,omitempty" json:"webhookAction,omitempty"`
}

func (c *Client) Actions() *Actions { return &Actions{client: c} }

func (n *Actions) List(viewName string) ([]Action, error) {
	var query struct {
		SearchDomain struct {
			Actions []Action `graphql:"actions"`
		} `graphql:"searchDomain(name: $viewName)"`
	}

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
	}

	err := n.client.Query(&query, variables)
	return query.SearchDomain.Actions, err
}

func (n *Actions) Update(viewName string, newAction *Action) (*Action, error) {
	if newAction == nil {
		return nil, fmt.Errorf("action must not be nil")
	}

	if newAction.ID == "" {
		return nil, fmt.Errorf("action must have non-empty action id")
	}

	currentAction, err := n.Get(viewName, newAction.Name)
	if err != nil {
		return nil, fmt.Errorf("unable to find action: %w", err)
	}

	if !reflect.ValueOf(newAction.EmailAction).IsZero() {
		var mutation struct {
			CreateEmailAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				EmailAction
			} `graphql:"updateEmailAction(input: { id: $id, viewName: $viewName, name: $actionName, recipients: $recipients, subjectTemplate: $subjectTemplate, bodyTemplate: $bodyTemplate, useProxy: $useProxy })"`
		}

		recipientsGQL := make([]graphql.String, len(newAction.EmailAction.Recipients))
		for i, recipient := range newAction.EmailAction.Recipients {
			recipientsGQL[i] = graphql.String(recipient)
		}
		variables := map[string]interface{}{
			"id":              graphql.String(currentAction.ID),
			"viewName":        graphql.String(viewName),
			"actionName":      graphql.String(newAction.Name),
			"recipients":      recipientsGQL,
			"subjectTemplate": graphql.String(newAction.EmailAction.SubjectTemplate),
			"bodyTemplate":    graphql.String(newAction.EmailAction.BodyTemplate),
			"useProxy":        graphql.Boolean(newAction.EmailAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateEmailAction.ID,
			Name: mutation.CreateEmailAction.Name,
			EmailAction: EmailAction{
				Recipients:      mutation.CreateEmailAction.Recipients,
				SubjectTemplate: mutation.CreateEmailAction.SubjectTemplate,
				BodyTemplate:    mutation.CreateEmailAction.BodyTemplate,
				UseProxy:        mutation.CreateEmailAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.HumioRepoAction).IsZero() {
		var mutation struct {
			CreateHumioRepoAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				HumioRepoAction
			} `graphql:"updateHumioRepoAction(input: { id: $id, viewName: $viewName, name: $actionName, ingestToken: $ingestToken })"`
		}

		variables := map[string]interface{}{
			"id":          graphql.String(currentAction.ID),
			"viewName":    graphql.String(viewName),
			"actionName":  graphql.String(newAction.Name),
			"ingestToken": graphql.String(newAction.HumioRepoAction.IngestToken),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateHumioRepoAction.ID,
			Name: mutation.CreateHumioRepoAction.Name,
			HumioRepoAction: HumioRepoAction{
				IngestToken: mutation.CreateHumioRepoAction.IngestToken,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.OpsGenieAction).IsZero() {
		var mutation struct {
			CreateOpsGenieAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				OpsGenieAction
			} `graphql:"updateOpsGenieAction(input: { id: $id, viewName: $viewName, name: $actionName, apiUrl: $apiUrl, genieKey: $genieKey, useProxy: $useProxy })"`
		}

		variables := map[string]interface{}{
			"id":         graphql.String(currentAction.ID),
			"viewName":   graphql.String(viewName),
			"actionName": graphql.String(newAction.Name),
			"apiUrl":     graphql.String(newAction.OpsGenieAction.ApiUrl),
			"genieKey":   graphql.String(newAction.OpsGenieAction.GenieKey),
			"useProxy":   graphql.Boolean(newAction.OpsGenieAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateOpsGenieAction.ID,
			Name: mutation.CreateOpsGenieAction.Name,
			OpsGenieAction: OpsGenieAction{
				ApiUrl:   mutation.CreateOpsGenieAction.ApiUrl,
				GenieKey: mutation.CreateOpsGenieAction.GenieKey,
				UseProxy: mutation.CreateOpsGenieAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.PagerDutyAction).IsZero() {
		var mutation struct {
			CreatePagerDutyAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				PagerDutyAction
			} `graphql:"updatePagerDutyAction(input: { id: $id, viewName: $viewName, name: $actionName, severity: $severity, routingKey: $routingKey, useProxy: $useProxy })"`
		}

		variables := map[string]interface{}{
			"id":         graphql.String(currentAction.ID),
			"viewName":   graphql.String(viewName),
			"actionName": graphql.String(newAction.Name),
			"severity":   graphql.String(newAction.PagerDutyAction.Severity),
			"routingKey": graphql.String(newAction.PagerDutyAction.RoutingKey),
			"useProxy":   graphql.Boolean(newAction.PagerDutyAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreatePagerDutyAction.ID,
			Name: mutation.CreatePagerDutyAction.Name,
			PagerDutyAction: PagerDutyAction{
				Severity:   mutation.CreatePagerDutyAction.Severity,
				RoutingKey: mutation.CreatePagerDutyAction.RoutingKey,
				UseProxy:   mutation.CreatePagerDutyAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.SlackAction).IsZero() {
		var mutation struct {
			CreateSlackAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				SlackAction
			} `graphql:"updateSlackAction(input: { id: $id, viewName: $viewName, name: $actionName, url: $url, fields: $fields, useProxy: $useProxy })"`
		}

		fieldsGQL := make([]SlackFieldEntryInput, len(newAction.SlackAction.Fields))
		for k, v := range newAction.SlackAction.Fields {
			fieldsGQL[k] = SlackFieldEntryInput{
				FieldName: v.FieldName,
				Value:     v.Value,
			}
		}
		variables := map[string]interface{}{
			"id":         graphql.String(currentAction.ID),
			"viewName":   graphql.String(viewName),
			"actionName": graphql.String(newAction.Name),
			"url":        graphql.String(newAction.SlackAction.Url),
			"fields":     fieldsGQL,
			"useProxy":   graphql.Boolean(newAction.SlackAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateSlackAction.ID,
			Name: mutation.CreateSlackAction.Name,
			SlackAction: SlackAction{
				Url:      mutation.CreateSlackAction.Url,
				Fields:   mutation.CreateSlackAction.Fields,
				UseProxy: mutation.CreateSlackAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.SlackPostMessageAction).IsZero() {
		var mutation struct {
			CreateSlackPostMessageAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				SlackPostMessageAction
			} `graphql:"updateSlackPostMessageAction(input: { id: $id, viewName: $viewName, name: $actionName, apiToken: $apiToken, channels: $channels, fields: $fields, useProxy: $useProxy })"`
		}

		channelsGQL := make([]graphql.String, len(newAction.SlackPostMessageAction.Channels))
		for k, v := range newAction.SlackPostMessageAction.Channels {
			channelsGQL[k] = graphql.String(v)
		}
		fieldsGQL := make([]SlackFieldEntryInput, len(newAction.SlackPostMessageAction.Fields))
		for k, v := range newAction.SlackPostMessageAction.Fields {
			fieldsGQL[k] = SlackFieldEntryInput{
				FieldName: v.FieldName,
				Value:     v.Value,
			}
		}
		variables := map[string]interface{}{
			"id":         graphql.String(currentAction.ID),
			"viewName":   graphql.String(viewName),
			"actionName": graphql.String(newAction.Name),
			"apiToken":   graphql.String(newAction.SlackPostMessageAction.ApiToken),
			"channels":   channelsGQL,
			"fields":     fieldsGQL,
			"useProxy":   graphql.Boolean(newAction.SlackPostMessageAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateSlackPostMessageAction.ID,
			Name: mutation.CreateSlackPostMessageAction.Name,
			SlackPostMessageAction: SlackPostMessageAction{
				ApiToken: mutation.CreateSlackPostMessageAction.ApiToken,
				Channels: mutation.CreateSlackPostMessageAction.Channels,
				Fields:   mutation.CreateSlackPostMessageAction.Fields,
				UseProxy: mutation.CreateSlackPostMessageAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.VictorOpsAction).IsZero() {
		var mutation struct {
			CreateVictorOpsAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				VictorOpsAction
			} `graphql:"updateVictorOpsAction(input: { id: $id, viewName: $viewName, name: $actionName, messageType: $messageType, notifyUrl: $notifyUrl, useProxy: $useProxy })"`
		}

		variables := map[string]interface{}{
			"id":          graphql.String(currentAction.ID),
			"viewName":    graphql.String(viewName),
			"actionName":  graphql.String(newAction.Name),
			"messageType": graphql.String(newAction.VictorOpsAction.MessageType),
			"notifyUrl":   graphql.String(newAction.VictorOpsAction.NotifyUrl),
			"useProxy":    graphql.Boolean(newAction.VictorOpsAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateVictorOpsAction.ID,
			Name: mutation.CreateVictorOpsAction.Name,
			VictorOpsAction: VictorOpsAction{
				MessageType: mutation.CreateVictorOpsAction.MessageType,
				NotifyUrl:   mutation.CreateVictorOpsAction.NotifyUrl,
				UseProxy:    mutation.CreateVictorOpsAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.WebhookAction).IsZero() {
		var mutation struct {
			CreateWebhookAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				WebhookAction
			} `graphql:"updateWebhookAction(input: { id: $id, viewName: $viewName, name: $actionName, url: $url, method: $method, headers: $headers, bodyTemplate: $bodyTemplate, ignoreSSL: $ignoreSSL, useProxy: $useProxy })"`
		}

		headersGQL := make([]HttpHeaderEntryInput, len(newAction.WebhookAction.Headers))
		for i, h := range newAction.WebhookAction.Headers {
			headersGQL[i] = HttpHeaderEntryInput{
				Header: h.Header,
				Value:  h.Value,
			}
		}
		variables := map[string]interface{}{
			"id":           graphql.String(currentAction.ID),
			"viewName":     graphql.String(viewName),
			"actionName":   graphql.String(newAction.Name),
			"url":          graphql.String(newAction.WebhookAction.Url),
			"method":       graphql.String(newAction.WebhookAction.Method),
			"headers":      headersGQL,
			"bodyTemplate": graphql.String(newAction.WebhookAction.BodyTemplate),
			"ignoreSSL":    graphql.Boolean(newAction.WebhookAction.IgnoreSSL),
			"useProxy":     graphql.Boolean(newAction.WebhookAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateWebhookAction.ID,
			Name: mutation.CreateWebhookAction.Name,
			WebhookAction: WebhookAction{
				Method:       mutation.CreateWebhookAction.Method,
				Url:          mutation.CreateWebhookAction.Url,
				Headers:      mutation.CreateWebhookAction.Headers,
				BodyTemplate: mutation.CreateWebhookAction.BodyTemplate,
				IgnoreSSL:    mutation.CreateWebhookAction.IgnoreSSL,
				UseProxy:     mutation.CreateWebhookAction.UseProxy,
			},
		}

		return &action, nil
	}

	return nil, fmt.Errorf("no action details specified or unsupported action type used")
}

func (n *Actions) Add(viewName string, newAction *Action) (*Action, error) {
	if newAction == nil {
		return nil, fmt.Errorf("action must not be nil")
	}

	if newAction.ID != "" {
		return nil, fmt.Errorf("action id must be empty when creating a new action")
	}

	if !reflect.ValueOf(newAction.EmailAction).IsZero() {
		var mutation struct {
			CreateEmailAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				EmailAction
			} `graphql:"createEmailAction(input: { viewName: $viewName, name: $actionName, recipients: $recipients, subjectTemplate: $subjectTemplate, bodyTemplate: $bodyTemplate, useProxy: $useProxy })"`
		}

		recipientsGQL := make([]graphql.String, len(newAction.EmailAction.Recipients))
		for i, recipient := range newAction.EmailAction.Recipients {
			recipientsGQL[i] = graphql.String(recipient)
		}
		variables := map[string]interface{}{
			"viewName":        graphql.String(viewName),
			"actionName":      graphql.String(newAction.Name),
			"recipients":      recipientsGQL,
			"subjectTemplate": graphql.String(newAction.EmailAction.SubjectTemplate),
			"bodyTemplate":    graphql.String(newAction.EmailAction.BodyTemplate),
			"useProxy":        graphql.Boolean(newAction.EmailAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateEmailAction.ID,
			Name: mutation.CreateEmailAction.Name,
			EmailAction: EmailAction{
				Recipients:      mutation.CreateEmailAction.Recipients,
				SubjectTemplate: mutation.CreateEmailAction.SubjectTemplate,
				BodyTemplate:    mutation.CreateEmailAction.BodyTemplate,
				UseProxy:        mutation.CreateEmailAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.HumioRepoAction).IsZero() {
		var mutation struct {
			CreateHumioRepoAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				HumioRepoAction
			} `graphql:"createHumioRepoAction(input: { viewName: $viewName, name: $actionName, ingestToken: $ingestToken })"`
		}

		variables := map[string]interface{}{
			"viewName":    graphql.String(viewName),
			"actionName":  graphql.String(newAction.Name),
			"ingestToken": graphql.String(newAction.HumioRepoAction.IngestToken),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateHumioRepoAction.ID,
			Name: mutation.CreateHumioRepoAction.Name,
			HumioRepoAction: HumioRepoAction{
				IngestToken: mutation.CreateHumioRepoAction.IngestToken,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.OpsGenieAction).IsZero() {
		var mutation struct {
			CreateOpsGenieAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				OpsGenieAction
			} `graphql:"createOpsGenieAction(input: { viewName: $viewName, name: $actionName, apiUrl: $apiUrl, genieKey: $genieKey, useProxy: $useProxy })"`
		}

		variables := map[string]interface{}{
			"viewName":   graphql.String(viewName),
			"actionName": graphql.String(newAction.Name),
			"apiUrl":     graphql.String(newAction.OpsGenieAction.ApiUrl),
			"genieKey":   graphql.String(newAction.OpsGenieAction.GenieKey),
			"useProxy":   graphql.Boolean(newAction.OpsGenieAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateOpsGenieAction.ID,
			Name: mutation.CreateOpsGenieAction.Name,
			OpsGenieAction: OpsGenieAction{
				ApiUrl:   mutation.CreateOpsGenieAction.ApiUrl,
				GenieKey: mutation.CreateOpsGenieAction.GenieKey,
				UseProxy: mutation.CreateOpsGenieAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.PagerDutyAction).IsZero() {
		var mutation struct {
			CreatePagerDutyAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				PagerDutyAction
			} `graphql:"createPagerDutyAction(input: { viewName: $viewName, name: $actionName, severity: $severity, routingKey: $routingKey, useProxy: $useProxy })"`
		}

		variables := map[string]interface{}{
			"viewName":   graphql.String(viewName),
			"actionName": graphql.String(newAction.Name),
			"severity":   graphql.String(newAction.PagerDutyAction.Severity),
			"routingKey": graphql.String(newAction.PagerDutyAction.RoutingKey),
			"useProxy":   graphql.Boolean(newAction.PagerDutyAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreatePagerDutyAction.ID,
			Name: mutation.CreatePagerDutyAction.Name,
			PagerDutyAction: PagerDutyAction{
				Severity:   mutation.CreatePagerDutyAction.Severity,
				RoutingKey: mutation.CreatePagerDutyAction.RoutingKey,
				UseProxy:   mutation.CreatePagerDutyAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.SlackAction).IsZero() {
		var mutation struct {
			CreateSlackAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				SlackAction
			} `graphql:"createSlackAction(input: { viewName: $viewName, name: $actionName, url: $url, fields: $fields, useProxy: $useProxy })"`
		}

		fieldsGQL := make([]SlackFieldEntryInput, len(newAction.SlackAction.Fields))
		for k, v := range newAction.SlackAction.Fields {
			fieldsGQL[k] = SlackFieldEntryInput{
				FieldName: v.FieldName,
				Value:     v.Value,
			}
		}
		variables := map[string]interface{}{
			"viewName":   graphql.String(viewName),
			"actionName": graphql.String(newAction.Name),
			"url":        graphql.String(newAction.SlackAction.Url),
			"fields":     fieldsGQL,
			"useProxy":   graphql.Boolean(newAction.SlackAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateSlackAction.ID,
			Name: mutation.CreateSlackAction.Name,
			SlackAction: SlackAction{
				Url:      mutation.CreateSlackAction.Url,
				Fields:   mutation.CreateSlackAction.Fields,
				UseProxy: mutation.CreateSlackAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.SlackPostMessageAction).IsZero() {
		var mutation struct {
			CreateSlackPostMessageAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				SlackPostMessageAction
			} `graphql:"createSlackPostMessageAction(input: { viewName: $viewName, name: $actionName, apiToken: $apiToken, channels: $channels, fields: $fields, useProxy: $useProxy })"`
		}

		channelsGQL := make([]graphql.String, len(newAction.SlackPostMessageAction.Channels))
		for k, v := range newAction.SlackPostMessageAction.Channels {
			channelsGQL[k] = graphql.String(v)
		}
		fieldsGQL := make([]SlackFieldEntryInput, len(newAction.SlackPostMessageAction.Fields))
		for k, v := range newAction.SlackPostMessageAction.Fields {
			fieldsGQL[k] = SlackFieldEntryInput{
				FieldName: v.FieldName,
				Value:     v.Value,
			}
		}
		variables := map[string]interface{}{
			"viewName":   graphql.String(viewName),
			"actionName": graphql.String(newAction.Name),
			"apiToken":   graphql.String(newAction.SlackPostMessageAction.ApiToken),
			"channels":   channelsGQL,
			"fields":     fieldsGQL,
			"useProxy":   graphql.Boolean(newAction.SlackPostMessageAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateSlackPostMessageAction.ID,
			Name: mutation.CreateSlackPostMessageAction.Name,
			SlackPostMessageAction: SlackPostMessageAction{
				ApiToken: mutation.CreateSlackPostMessageAction.ApiToken,
				Channels: mutation.CreateSlackPostMessageAction.Channels,
				Fields:   mutation.CreateSlackPostMessageAction.Fields,
				UseProxy: mutation.CreateSlackPostMessageAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.VictorOpsAction).IsZero() {
		var mutation struct {
			CreateVictorOpsAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				VictorOpsAction
			} `graphql:"createVictorOpsAction(input: { viewName: $viewName, name: $actionName, messageType: $messageType, notifyUrl: $notifyUrl, useProxy: $useProxy })"`
		}

		variables := map[string]interface{}{
			"viewName":    graphql.String(viewName),
			"actionName":  graphql.String(newAction.Name),
			"messageType": graphql.String(newAction.VictorOpsAction.MessageType),
			"notifyUrl":   graphql.String(newAction.VictorOpsAction.NotifyUrl),
			"useProxy":    graphql.Boolean(newAction.VictorOpsAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateVictorOpsAction.ID,
			Name: mutation.CreateVictorOpsAction.Name,
			VictorOpsAction: VictorOpsAction{
				MessageType: mutation.CreateVictorOpsAction.MessageType,
				NotifyUrl:   mutation.CreateVictorOpsAction.NotifyUrl,
				UseProxy:    mutation.CreateVictorOpsAction.UseProxy,
			},
		}

		return &action, nil
	}

	if !reflect.ValueOf(newAction.WebhookAction).IsZero() {
		var mutation struct {
			CreateWebhookAction struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				WebhookAction
			} `graphql:"createWebhookAction(input: { viewName: $viewName, name: $actionName, url: $url, method: $method, headers: $headers, bodyTemplate: $bodyTemplate, ignoreSSL: $ignoreSSL, useProxy: $useProxy })"`
		}

		headersGQL := make([]HttpHeaderEntryInput, len(newAction.WebhookAction.Headers))
		for i, h := range newAction.WebhookAction.Headers {
			headersGQL[i] = HttpHeaderEntryInput{
				Header: h.Header,
				Value:  h.Value,
			}
		}
		variables := map[string]interface{}{
			"viewName":     graphql.String(viewName),
			"actionName":   graphql.String(newAction.Name),
			"url":          graphql.String(newAction.WebhookAction.Url),
			"method":       graphql.String(newAction.WebhookAction.Method),
			"headers":      headersGQL,
			"bodyTemplate": graphql.String(newAction.WebhookAction.BodyTemplate),
			"ignoreSSL":    graphql.Boolean(newAction.WebhookAction.IgnoreSSL),
			"useProxy":     graphql.Boolean(newAction.WebhookAction.UseProxy),
		}

		err := n.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		action := Action{
			ID:   mutation.CreateWebhookAction.ID,
			Name: mutation.CreateWebhookAction.Name,
			WebhookAction: WebhookAction{
				Method:       mutation.CreateWebhookAction.Method,
				Url:          mutation.CreateWebhookAction.Url,
				Headers:      mutation.CreateWebhookAction.Headers,
				BodyTemplate: mutation.CreateWebhookAction.BodyTemplate,
				IgnoreSSL:    mutation.CreateWebhookAction.IgnoreSSL,
				UseProxy:     mutation.CreateWebhookAction.UseProxy,
			},
		}

		return &action, nil
	}

	return nil, fmt.Errorf("no action details specified or unsupported action type used")
}

func (n *Actions) Get(viewName, actionName string) (*Action, error) {
	actions, err := n.List(viewName)
	if err != nil {
		return nil, fmt.Errorf("unable to list actions: %w", err)
	}
	for _, action := range actions {
		if action.Name == actionName {
			return &action, nil
		}
	}

	return nil, ActionNotFound(actionName)
}

func (n *Actions) GetByID(viewName, actionID string) (*Action, error) {
	var query struct {
		SearchDomain struct {
			Action *Action `graphql:"action(id: $actionId)"`
		} `graphql:"searchDomain(name: $viewName)"`
	}

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
		"actionId": graphql.String(actionID),
	}

	err := n.client.Query(&query, variables)

	if err != nil {
		return nil, err
	}

	if query.SearchDomain.Action == nil {
		return nil, ActionNotFound(actionID)
	}
	return query.SearchDomain.Action, err
}

func (n *Actions) Delete(viewName, actionName string) error {
	actions, err := n.List(viewName)
	if err != nil {
		return fmt.Errorf("unable to list actions: %w", err)
	}
	var actionID string
	for _, action := range actions {
		if action.Name == actionName {
			actionID = action.ID
			break
		}
	}
	if actionID == "" {
		return fmt.Errorf("unable to find action")
	}

	var mutation struct {
		DeleteAction bool `graphql:"deleteAction(input: { viewName: $viewName, id: $actionId })"`
	}

	variables := map[string]interface{}{
		"viewName": graphql.String(viewName),
		"actionId": graphql.String(actionID),
	}

	return n.client.Mutate(&mutation, variables)
}
