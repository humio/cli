package humiographql

import graphql "github.com/cli/shurcooL-graphql"

type ScheduledSearch struct {
	ID                         graphql.String   `json:"id"`
	Name                       graphql.String   `json:"name"`
	Description                graphql.String   `json:"description,omitempty"`
	QueryString                graphql.String   `json:"queryString"`
	Start                      graphql.String   `json:"start"`
	End                        graphql.String   `json:"end"`
	TimeZone                   graphql.String   `json:"timezone"`
	Schedule                   graphql.String   `json:"schedule"`
	BackfillLimit              graphql.Int      `json:"backfillLimit"`
	Enabled                    graphql.Boolean  `json:"enabled"`
	ActionsV2                  []Action         `json:"actionsV2"`
	RunAsUser                  User             `json:"runAsUser,omitempty"`
	TimeOfNextPlannedExecution Long             `json:"timeOfNextPlannedExecution"`
	Labels                     []graphql.String `json:"labels"`
	QueryOwnership             QueryOwnership   `json:"queryOwnership"`
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
