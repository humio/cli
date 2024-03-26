package humiographql

import graphql "github.com/cli/shurcooL-graphql"

type AggregateAlert struct {
	ID                    graphql.String     `graphql:"id"`
	Name                  graphql.String     `graphql:"name"`
	Description           graphql.String     `graphql:"description"`
	QueryString           graphql.String     `graphql:"queryString"`
	SearchIntervalSeconds Long               `graphql:"searchIntervalSeconds"`
	ThrottleTimeSeconds   Long               `graphql:"throttleTimeSeconds"`
	ThrottleField         graphql.String     `graphql:"throttleField"`
	Actions               []Action           `graphql:"actions"`
	Labels                []graphql.String   `graphql:"labels"`
	Enabled               graphql.Boolean    `graphql:"enabled"`
	QueryOwnership        QueryOwnership     `graphql:"queryOwnership"`
	TriggerMode           TriggerMode        `graphql:"triggerMode"`
	QueryTimestampType    QueryTimestampType `graphql:"queryTimestampType"`
}

type CreateAggregateAlert struct {
	ViewName              RepoOrViewName     `json:"viewName"`
	Name                  graphql.String     `json:"name"`
	Description           graphql.String     `json:"description,omitempty"`
	QueryString           graphql.String     `json:"queryString"`
	SearchIntervalSeconds Long               `json:"searchIntervalSeconds"`
	ThrottleTimeSeconds   Long               `json:"throttleTimeSeconds"`
	ThrottleField         graphql.String     `json:"throttleField,omitempty"`
	ActionIdsOrNames      []graphql.String   `json:"actionIdsOrNames"`
	Labels                []graphql.String   `json:"labels"`
	Enabled               graphql.Boolean    `json:"enabled"`
	RunAsUserID           graphql.String     `json:"runAsUserId,omitempty"`
	QueryOwnershipType    QueryOwnershipType `json:"queryOwnershipType"`
	TriggerMode           TriggerMode        `json:"triggerMode,omitempty"`
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
	ThrottleField         graphql.String     `json:"throttleField,omitempty"`
	ActionIdsOrNames      []graphql.String   `json:"actionIdsOrNames"`
	Labels                []graphql.String   `json:"labels"`
	Enabled               graphql.Boolean    `json:"enabled"`
	RunAsUserID           graphql.String     `json:"runAsUserId,omitempty"`
	QueryOwnershipType    QueryOwnershipType `json:"queryOwnershipType"`
	TriggerMode           TriggerMode        `json:"triggerMode"`
	QueryTimestampType    QueryTimestampType `json:"queryTimestampType"`
}
