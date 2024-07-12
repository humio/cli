package humiographql

import graphql "github.com/cli/shurcooL-graphql"

type ScheduledSearch struct {
	ID             graphql.String   `graphql:"id"`
	Name           graphql.String   `graphql:"name"`
	Description    graphql.String   `graphql:"description"`
	QueryString    graphql.String   `graphql:"queryString"`
	Start          graphql.String   `graphql:"start"`
	End            graphql.String   `graphql:"end"`
	TimeZone       graphql.String   `graphql:"timeZone"`
	Schedule       graphql.String   `graphql:"schedule"`
	BackfillLimit  graphql.Int      `graphql:"backfillLimit"`
	Enabled        graphql.Boolean  `graphql:"enabled"`
	ActionsV2      []Action         `graphql:"actionsV2"`
	RunAsUser      User             `graphql:"runAsUser"`
	Labels         []graphql.String `graphql:"labels"`
	QueryOwnership QueryOwnership   `graphql:"queryOwnership"`
}

type CreateScheduledSearch struct {
	ViewName          graphql.String     `json:"viewName"`
	Name              graphql.String     `json:"name"`
	Description       graphql.String     `json:"description,omitempty"`
	QueryString       graphql.String     `json:"queryString"`
	QueryStart        graphql.String     `json:"queryStart"`
	QueryEnd          graphql.String     `json:"queryEnd"`
	Schedule          graphql.String     `json:"schedule"`
	TimeZone          graphql.String     `json:"timeZone"`
	BackfillLimit     graphql.Int        `json:"backfillLimit"`
	Enabled           graphql.Boolean    `json:"enabled"`
	ActionsIdsOrNames []graphql.String   `json:"actions"`
	Labels            []graphql.String   `json:"labels,omitempty"`
	RunAsUserID       graphql.String     `json:"runAsUserId,omitempty"`
	QueryOwnership    QueryOwnershipType `json:"queryOwnershipType,omitempty"`
}

type UpdateScheduledSearch struct {
	ViewName          graphql.String     `json:"viewName"`
	ID                graphql.String     `json:"id"`
	Name              graphql.String     `json:"name"`
	Description       graphql.String     `json:"description,omitempty"`
	QueryString       graphql.String     `json:"queryString"`
	QueryStart        graphql.String     `json:"queryStart"`
	QueryEnd          graphql.String     `json:"queryEnd"`
	Schedule          graphql.String     `json:"schedule"`
	TimeZone          graphql.String     `json:"timeZone"`
	BackfillLimit     graphql.Int        `json:"backfillLimit"`
	Enabled           graphql.Boolean    `json:"enabled"`
	ActionsIdsOrNames []graphql.String   `json:"actions"`
	Labels            []graphql.String   `json:"labels"`
	RunAsUserID       graphql.String     `json:"runAsUserId,omitempty"`
	QueryOwnership    QueryOwnershipType `json:"queryOwnershipType"`
}
