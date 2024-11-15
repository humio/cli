package api

import (
	"context"
	"fmt"
	"reflect"

	"github.com/humio/cli/internal/api/humiographql"
)

type Actions struct {
	client *Client
}

type EmailAction struct {
	Recipients      []string
	SubjectTemplate *string
	BodyTemplate    *string
	UseProxy        bool
}

type HumioRepoAction struct {
	IngestToken string
}

type OpsGenieAction struct {
	ApiUrl   string
	GenieKey string
	UseProxy bool
}

type PagerDutyAction struct {
	Severity   string
	RoutingKey string
	UseProxy   bool
}

type SlackField struct {
	FieldName string
	Value     string
}

type SlackAction struct {
	Url      string
	Fields   []SlackField
	UseProxy bool
}

type SlackPostMessageAction struct {
	ApiToken string
	Channels []string
	Fields   []SlackField
	UseProxy bool
}

type UploadFileAction struct {
	FileName string
}

type VictorOpsAction struct {
	MessageType string
	NotifyUrl   string
	UseProxy    bool
}

type HttpHeader struct {
	Header string
	Value  string
}

type WebhookAction struct {
	Method       string
	Url          string
	Headers      []HttpHeader
	BodyTemplate string
	IgnoreSSL    bool
	UseProxy     bool
}

type Action struct {
	Type string
	ID   string `yaml:"-"`
	Name string

	EmailAction            EmailAction            `yaml:"emailAction,omitempty"`
	HumioRepoAction        HumioRepoAction        `yaml:"humioRepoAction,omitempty"`
	OpsGenieAction         OpsGenieAction         `yaml:"opsGenieAction,omitempty"`
	PagerDutyAction        PagerDutyAction        `yaml:"pagerDutyAction,omitempty"`
	SlackAction            SlackAction            `yaml:"slackAction,omitempty"`
	SlackPostMessageAction SlackPostMessageAction `yaml:"slackPostMessageAction,omitempty"`
	VictorOpsAction        VictorOpsAction        `yaml:"victorOpsAction,omitempty"`
	UploadFileAction       UploadFileAction       `yaml:"uploadFileAction,omitempty"`
	WebhookAction          WebhookAction          `yaml:"webhookAction,omitempty"`
}

func (c *Client) Actions() *Actions { return &Actions{client: c} }

func (n *Actions) List(searchDomainName string) ([]Action, error) {
	resp, err := humiographql.ListActions(context.Background(), n.client, searchDomainName)
	if err != nil {
		return nil, err
	}
	respSearchDomain := resp.GetSearchDomain()
	respSearchDomainActions := respSearchDomain.GetActions()
	actions := make([]Action, len(respSearchDomainActions))
	for idx, action := range respSearchDomainActions {
		switch v := action.(type) {
		case *humiographql.ListActionsSearchDomainActionsEmailAction:
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
				EmailAction: EmailAction{
					Recipients:      v.GetRecipients(),
					SubjectTemplate: v.GetSubjectTemplate(),
					BodyTemplate:    v.GetEmailBodyTemplate(),
					UseProxy:        v.GetUseProxy(),
				},
			}
		case *humiographql.ListActionsSearchDomainActionsHumioRepoAction:
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
				HumioRepoAction: HumioRepoAction{
					IngestToken: v.GetIngestToken(),
				},
			}
		case *humiographql.ListActionsSearchDomainActionsOpsGenieAction:
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
				OpsGenieAction: OpsGenieAction{
					ApiUrl:   v.GetApiUrl(),
					GenieKey: v.GetGenieKey(),
					UseProxy: v.GetUseProxy(),
				},
			}
		case *humiographql.ListActionsSearchDomainActionsPagerDutyAction:
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
				PagerDutyAction: PagerDutyAction{
					Severity:   v.GetSeverity(),
					RoutingKey: v.GetRoutingKey(),
					UseProxy:   v.GetUseProxy(),
				},
			}
		case *humiographql.ListActionsSearchDomainActionsSlackAction:
			fields := make([]SlackField, len(v.GetFields()))
			for jdx, field := range v.GetFields() {
				fields[jdx] = SlackField{
					FieldName: field.GetFieldName(),
					Value:     field.GetValue(),
				}
			}
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
				SlackAction: SlackAction{
					Url:      v.GetUrl(),
					Fields:   fields,
					UseProxy: v.GetUseProxy(),
				},
			}
		case *humiographql.ListActionsSearchDomainActionsSlackPostMessageAction:
			fields := make([]SlackField, len(v.GetFields()))
			for jdx, field := range v.GetFields() {
				fields[jdx] = SlackField{
					FieldName: field.GetFieldName(),
					Value:     field.GetValue(),
				}
			}
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
				SlackPostMessageAction: SlackPostMessageAction{
					ApiToken: v.GetApiToken(),
					Channels: v.GetChannels(),
					Fields:   fields,
					UseProxy: v.GetUseProxy(),
				},
			}
		case *humiographql.ListActionsSearchDomainActionsVictorOpsAction:
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
				VictorOpsAction: VictorOpsAction{
					MessageType: v.GetMessageType(),
					NotifyUrl:   v.GetNotifyUrl(),
					UseProxy:    v.GetUseProxy(),
				},
			}
		case *humiographql.ListActionsSearchDomainActionsUploadFileAction:
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
				UploadFileAction: UploadFileAction{
					FileName: v.GetFileName(),
				},
			}
		case *humiographql.ListActionsSearchDomainActionsWebhookAction:
			headers := make([]HttpHeader, len(v.GetHeaders()))
			for jdx, header := range v.GetHeaders() {
				headers[jdx] = HttpHeader{
					Header: header.GetHeader(),
					Value:  header.GetValue(),
				}
			}
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
				WebhookAction: WebhookAction{
					Method:       v.GetMethod(),
					Url:          v.GetUrl(),
					Headers:      headers,
					BodyTemplate: v.GetWebhookBodyTemplate(),
					IgnoreSSL:    v.GetIgnoreSSL(),
					UseProxy:     v.GetUseProxy(),
				},
			}
		default:
			actions[idx] = Action{
				Type: *v.GetTypename(),
				ID:   v.GetId(),
				Name: v.GetName(),
			}
		}
	}

	return actions, nil
}

func (n *Actions) Update(searchDomainName string, newAction *Action) (*Action, error) {
	if newAction == nil {
		return nil, fmt.Errorf("action must not be nil")
	}

	if newAction.ID == "" {
		return nil, fmt.Errorf("action must have non-empty action id")
	}

	currentAction, getErr := n.Get(searchDomainName, newAction.Name)
	if getErr != nil {
		return nil, ActionNotFound(newAction.Name)
	}

	if !reflect.ValueOf(newAction.EmailAction).IsZero() {
		resp, err := humiographql.UpdateEmailAction(
			context.Background(),
			n.client,
			searchDomainName,
			currentAction.ID,
			newAction.Name,
			newAction.EmailAction.Recipients,
			newAction.EmailAction.SubjectTemplate,
			newAction.EmailAction.BodyTemplate,
			newAction.EmailAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetUpdateEmailAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			EmailAction: EmailAction{
				Recipients:      respUpdate.GetRecipients(),
				SubjectTemplate: respUpdate.GetSubjectTemplate(),
				BodyTemplate:    respUpdate.GetBodyTemplate(),
				UseProxy:        respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.HumioRepoAction).IsZero() {
		resp, err := humiographql.UpdateHumioRepoAction(
			context.Background(),
			n.client,
			searchDomainName,
			currentAction.ID,
			newAction.Name,
			newAction.HumioRepoAction.IngestToken,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetUpdateHumioRepoAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			HumioRepoAction: HumioRepoAction{
				IngestToken: respUpdate.GetIngestToken(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.OpsGenieAction).IsZero() {
		resp, err := humiographql.UpdateOpsGenieAction(
			context.Background(),
			n.client,
			searchDomainName,
			currentAction.ID,
			newAction.Name,
			newAction.OpsGenieAction.ApiUrl,
			newAction.OpsGenieAction.GenieKey,
			newAction.OpsGenieAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetUpdateOpsGenieAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			OpsGenieAction: OpsGenieAction{
				ApiUrl:   respUpdate.GetApiUrl(),
				GenieKey: respUpdate.GetGenieKey(),
				UseProxy: respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.PagerDutyAction).IsZero() {
		resp, err := humiographql.UpdatePagerDutyAction(
			context.Background(),
			n.client,
			searchDomainName,
			currentAction.ID,
			newAction.Name,
			newAction.PagerDutyAction.Severity,
			newAction.PagerDutyAction.RoutingKey,
			newAction.PagerDutyAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetUpdatePagerDutyAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			PagerDutyAction: PagerDutyAction{
				Severity:   respUpdate.GetSeverity(),
				RoutingKey: respUpdate.GetRoutingKey(),
				UseProxy:   respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.SlackAction).IsZero() {
		fields := make([]humiographql.SlackFieldEntryInput, len(newAction.SlackAction.Fields))
		for idx, field := range newAction.SlackAction.Fields {
			fields[idx] = humiographql.SlackFieldEntryInput{
				FieldName: field.FieldName,
				Value:     field.Value,
			}
		}
		resp, err := humiographql.UpdateSlackAction(
			context.Background(),
			n.client,
			searchDomainName,
			currentAction.ID,
			newAction.Name,
			fields,
			newAction.SlackAction.Url,
			newAction.SlackAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetUpdateSlackAction()
		respUpdateFields := respUpdate.GetFields()
		fieldsUpdate := make([]SlackField, len(respUpdateFields))
		for idx, field := range respUpdateFields {
			fieldsUpdate[idx] = SlackField{
				FieldName: field.GetFieldName(),
				Value:     field.GetValue(),
			}
		}
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			SlackAction: SlackAction{
				Fields:   fieldsUpdate,
				Url:      respUpdate.GetUrl(),
				UseProxy: respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.SlackPostMessageAction).IsZero() {
		fields := make([]humiographql.SlackFieldEntryInput, len(newAction.SlackPostMessageAction.Fields))
		for idx, field := range newAction.SlackPostMessageAction.Fields {
			fields[idx] = humiographql.SlackFieldEntryInput{
				FieldName: field.FieldName,
				Value:     field.Value,
			}
		}
		resp, err := humiographql.UpdateSlackPostMessageAction(
			context.Background(),
			n.client,
			searchDomainName,
			currentAction.ID,
			newAction.Name,
			newAction.SlackPostMessageAction.ApiToken,
			newAction.SlackPostMessageAction.Channels,
			fields,
			newAction.SlackPostMessageAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetUpdateSlackPostMessageAction()
		respUpdateFields := respUpdate.GetFields()
		fieldsUpdate := make([]SlackField, len(respUpdateFields))
		for idx, field := range respUpdateFields {
			fieldsUpdate[idx] = SlackField{
				FieldName: field.GetFieldName(),
				Value:     field.GetValue(),
			}
		}
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			SlackPostMessageAction: SlackPostMessageAction{
				ApiToken: respUpdate.GetApiToken(),
				Channels: respUpdate.GetChannels(),
				Fields:   fieldsUpdate,
				UseProxy: respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.VictorOpsAction).IsZero() {
		resp, err := humiographql.UpdateVictorOpsAction(
			context.Background(),
			n.client,
			searchDomainName,
			currentAction.ID,
			newAction.Name,
			newAction.VictorOpsAction.MessageType,
			newAction.VictorOpsAction.NotifyUrl,
			newAction.VictorOpsAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetUpdateVictorOpsAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			VictorOpsAction: VictorOpsAction{
				MessageType: respUpdate.GetMessageType(),
				NotifyUrl:   respUpdate.GetNotifyUrl(),
				UseProxy:    respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.UploadFileAction).IsZero() {
		resp, err := humiographql.UpdateUploadFileAction(
			context.Background(),
			n.client,
			searchDomainName,
			currentAction.ID,
			newAction.Name,
			newAction.UploadFileAction.FileName,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetUpdateUploadFileAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			UploadFileAction: UploadFileAction{
				FileName: respUpdate.GetFileName(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.WebhookAction).IsZero() {
		headers := make([]humiographql.HttpHeaderEntryInput, len(newAction.WebhookAction.Headers))
		for idx, header := range newAction.WebhookAction.Headers {
			headers[idx] = humiographql.HttpHeaderEntryInput{
				Header: header.Header,
				Value:  header.Value,
			}
		}
		resp, err := humiographql.UpdateWebhookAction(
			context.Background(),
			n.client,
			searchDomainName,
			currentAction.ID,
			newAction.Name,
			newAction.WebhookAction.Url,
			newAction.WebhookAction.Method,
			headers,
			newAction.WebhookAction.BodyTemplate,
			newAction.WebhookAction.IgnoreSSL,
			newAction.WebhookAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetUpdateWebhookAction()
		respUpdateHeaders := respUpdate.GetHeaders()
		fieldsUpdate := make([]HttpHeader, len(respUpdateHeaders))
		for idx, header := range respUpdateHeaders {
			fieldsUpdate[idx] = HttpHeader{
				Header: header.GetHeader(),
				Value:  header.GetValue(),
			}
		}
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			WebhookAction: WebhookAction{
				Url:          respUpdate.GetUrl(),
				Method:       respUpdate.GetMethod(),
				Headers:      fieldsUpdate,
				BodyTemplate: respUpdate.GetBodyTemplate(),
				IgnoreSSL:    respUpdate.GetIgnoreSSL(),
				UseProxy:     respUpdate.GetUseProxy(),
			},
		}, nil
	}

	return nil, fmt.Errorf("no action details specified or unsupported action type used")
}

func (n *Actions) Add(searchDomainName string, newAction *Action) (*Action, error) {
	if newAction == nil {
		return nil, fmt.Errorf("action must not be nil")
	}

	if !reflect.ValueOf(newAction.EmailAction).IsZero() {
		resp, err := humiographql.CreateEmailAction(
			context.Background(),
			n.client,
			searchDomainName,
			newAction.Name,
			newAction.EmailAction.Recipients,
			newAction.EmailAction.SubjectTemplate,
			newAction.EmailAction.BodyTemplate,
			newAction.EmailAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetCreateEmailAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			EmailAction: EmailAction{
				Recipients:      respUpdate.GetRecipients(),
				SubjectTemplate: respUpdate.GetSubjectTemplate(),
				BodyTemplate:    respUpdate.GetBodyTemplate(),
				UseProxy:        respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.HumioRepoAction).IsZero() {
		resp, err := humiographql.CreateHumioRepoAction(
			context.Background(),
			n.client,
			searchDomainName,
			newAction.Name,
			newAction.HumioRepoAction.IngestToken,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetCreateHumioRepoAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			HumioRepoAction: HumioRepoAction{
				IngestToken: respUpdate.GetIngestToken(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.OpsGenieAction).IsZero() {
		resp, err := humiographql.CreateOpsGenieAction(
			context.Background(),
			n.client,
			searchDomainName,
			newAction.Name,
			newAction.OpsGenieAction.ApiUrl,
			newAction.OpsGenieAction.GenieKey,
			newAction.OpsGenieAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetCreateOpsGenieAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			OpsGenieAction: OpsGenieAction{
				ApiUrl:   respUpdate.GetApiUrl(),
				GenieKey: respUpdate.GetGenieKey(),
				UseProxy: respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.PagerDutyAction).IsZero() {
		resp, err := humiographql.CreatePagerDutyAction(
			context.Background(),
			n.client,
			searchDomainName,
			newAction.Name,
			newAction.PagerDutyAction.Severity,
			newAction.PagerDutyAction.RoutingKey,
			newAction.PagerDutyAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetCreatePagerDutyAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			PagerDutyAction: PagerDutyAction{
				Severity:   respUpdate.GetSeverity(),
				RoutingKey: respUpdate.GetRoutingKey(),
				UseProxy:   respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.SlackAction).IsZero() {
		fields := make([]humiographql.SlackFieldEntryInput, len(newAction.SlackAction.Fields))
		for idx, field := range newAction.SlackAction.Fields {
			fields[idx] = humiographql.SlackFieldEntryInput{
				FieldName: field.FieldName,
				Value:     field.Value,
			}
		}
		resp, err := humiographql.CreateSlackAction(
			context.Background(),
			n.client,
			searchDomainName,
			newAction.Name,
			fields,
			newAction.SlackAction.Url,
			newAction.SlackAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetCreateSlackAction()
		respUpdateFields := respUpdate.GetFields()
		fieldsUpdate := make([]SlackField, len(respUpdateFields))
		for idx, field := range respUpdateFields {
			fieldsUpdate[idx] = SlackField{
				FieldName: field.GetFieldName(),
				Value:     field.GetValue(),
			}
		}
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			SlackAction: SlackAction{
				Fields:   fieldsUpdate,
				Url:      respUpdate.GetUrl(),
				UseProxy: respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.SlackPostMessageAction).IsZero() {
		fields := make([]humiographql.SlackFieldEntryInput, len(newAction.SlackPostMessageAction.Fields))
		for idx, field := range newAction.SlackPostMessageAction.Fields {
			fields[idx] = humiographql.SlackFieldEntryInput{
				FieldName: field.FieldName,
				Value:     field.Value,
			}
		}
		resp, err := humiographql.CreateSlackPostMessageAction(
			context.Background(),
			n.client,
			searchDomainName,
			newAction.Name,
			newAction.SlackPostMessageAction.ApiToken,
			newAction.SlackPostMessageAction.Channels,
			fields,
			newAction.SlackPostMessageAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetCreateSlackPostMessageAction()
		respUpdateFields := respUpdate.GetFields()
		fieldsUpdate := make([]SlackField, len(respUpdateFields))
		for idx, field := range respUpdateFields {
			fieldsUpdate[idx] = SlackField{
				FieldName: field.GetFieldName(),
				Value:     field.GetValue(),
			}
		}
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			SlackPostMessageAction: SlackPostMessageAction{
				ApiToken: respUpdate.GetApiToken(),
				Channels: respUpdate.GetChannels(),
				Fields:   fieldsUpdate,
				UseProxy: respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.VictorOpsAction).IsZero() {
		resp, err := humiographql.CreateVictorOpsAction(
			context.Background(),
			n.client,
			searchDomainName,
			newAction.Name,
			newAction.VictorOpsAction.MessageType,
			newAction.VictorOpsAction.NotifyUrl,
			newAction.VictorOpsAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetCreateVictorOpsAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			VictorOpsAction: VictorOpsAction{
				MessageType: respUpdate.GetMessageType(),
				NotifyUrl:   respUpdate.GetNotifyUrl(),
				UseProxy:    respUpdate.GetUseProxy(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.UploadFileAction).IsZero() {
		resp, err := humiographql.CreateUploadFileAction(
			context.Background(),
			n.client,
			searchDomainName,
			newAction.Name,
			newAction.UploadFileAction.FileName,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetCreateUploadFileAction()
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			UploadFileAction: UploadFileAction{
				FileName: respUpdate.GetFileName(),
			},
		}, nil
	}

	if !reflect.ValueOf(newAction.WebhookAction).IsZero() {
		headers := make([]humiographql.HttpHeaderEntryInput, len(newAction.WebhookAction.Headers))
		for idx, header := range newAction.WebhookAction.Headers {
			headers[idx] = humiographql.HttpHeaderEntryInput{
				Header: header.Header,
				Value:  header.Value,
			}
		}
		resp, err := humiographql.CreateWebhookAction(
			context.Background(),
			n.client,
			searchDomainName,
			newAction.Name,
			newAction.WebhookAction.Url,
			newAction.WebhookAction.Method,
			headers,
			newAction.WebhookAction.BodyTemplate,
			newAction.WebhookAction.IgnoreSSL,
			newAction.WebhookAction.UseProxy,
		)
		if err != nil {
			return nil, err
		}

		respUpdate := resp.GetCreateWebhookAction()
		respUpdateHeaders := respUpdate.GetHeaders()
		fieldsUpdate := make([]HttpHeader, len(respUpdateHeaders))
		for idx, header := range respUpdateHeaders {
			fieldsUpdate[idx] = HttpHeader{
				Header: header.GetHeader(),
				Value:  header.GetValue(),
			}
		}
		return &Action{
			ID:   respUpdate.GetId(),
			Name: respUpdate.GetName(),
			WebhookAction: WebhookAction{
				Url:          respUpdate.GetUrl(),
				Method:       respUpdate.GetMethod(),
				Headers:      fieldsUpdate,
				BodyTemplate: respUpdate.GetBodyTemplate(),
				IgnoreSSL:    respUpdate.GetIgnoreSSL(),
				UseProxy:     respUpdate.GetUseProxy(),
			},
		}, nil
	}

	return nil, fmt.Errorf("no action details specified or unsupported action type used")
}

func (n *Actions) Get(searchDomainName, actionName string) (*Action, error) {
	actions, err := n.List(searchDomainName)
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

func (n *Actions) Delete(searchDomainName, actionName string) error {
	actions, err := n.List(searchDomainName)
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
		return ActionNotFound(actionID)
	}

	_, err = humiographql.DeleteActionByID(context.Background(), n.client, searchDomainName, actionID)
	if err != nil {
		return err
	}

	return nil
}
