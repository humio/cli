package humiographql

import graphql "github.com/cli/shurcooL-graphql"

type AggregateAlert struct {
	ID                    graphql.String     `json:"id"`
	Name                  graphql.String     `json:"name"`
	Description           graphql.String     `json:"description,omitempty"`
	QueryString           graphql.String     `json:"queryString"`
	SearchIntervalSeconds Long               `json:"searchIntervalSeconds"`
	ThrottleTimeSeconds   Long               `json:"throttleTimeSeconds"`
	ThrottleField         graphql.String     `json:"throttleField"`
	Actions               []Action           `json:"actionIdsOrNames"`
	Labels                []graphql.String   `json:"labels"`
	Enabled               graphql.Boolean    `json:"enabled"`
	QueryOwnership        QueryOwnership     `json:"queryOwnershipType"`
	TriggerMode           TriggerMode        `json:"triggerMode"`
	QueryTimestampType    QueryTimestampType `json:"queryTimestampType"`
}

type CreateAggregateAlert struct {
	ViewName              RepoOrViewName     `json:"viewName"`
	Name                  graphql.String     `json:"name"`
	Description           graphql.String     `json:"description,omitempty"`
	QueryString           graphql.String     `json:"queryString"`
	SearchIntervalSeconds Long               `json:"searchIntervalSeconds"`
	ThrottleTimeSeconds   Long               `json:"throttleTimeSeconds"`
	ThrottleField         graphql.String     `json:"throttleField"`
	ActionIdsOrNames      []graphql.String   `json:"actionIdsOrNames"`
	Labels                []graphql.String   `json:"labels"`
	Enabled               graphql.Boolean    `json:"enabled"`
	RunAsUserID           graphql.String     `json:"runAsUserId,omitempty"`
	QueryOwnershipType    QueryOwnershipType `json:"queryOwnershipType"`
	TriggerMode           TriggerMode        `json:"triggerMode"`
	QueryTimestampType    QueryTimestampType `json:"queryTimestampType"`
}

type UpdateAggregateAlert struct {
	ViewName              RepoOrViewName     `json:"viewName"`
	ID                    graphql.String     `json:"id"`
	Name                  graphql.String     `json:"name"`
	Description           graphql.String     `json:"description,omitempty"`
	QueryString           graphql.String     `json:"queryString"`
	SearchIntervalSeconds Long               `json:"searchIntervalSeconds"`
	ThrottleTimeSeconds   Long               `json:"throttleTimeSeconds"`
	ThrottleField         graphql.String     `json:"throttleField"`
	ActionIdsOrNames      []graphql.String   `json:"actionIdsOrNames"`
	Labels                []graphql.String   `json:"labels"`
	Enabled               graphql.Boolean    `json:"enabled"`
	RunAsUserID           graphql.String     `json:"runAsUserId,omitempty"`
	QueryOwnershipType    QueryOwnershipType `json:"queryOwnershipType"`
	TriggerMode           TriggerMode        `json:"triggerMode"`
	QueryTimestampType    QueryTimestampType `json:"queryTimestampType"`
}
