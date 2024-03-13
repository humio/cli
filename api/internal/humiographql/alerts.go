package humiographql

import graphql "github.com/cli/shurcooL-graphql"

type Alert struct {
	ID                 graphql.String   `graphql:"id"`
	Name               graphql.String   `graphql:"name"`
	QueryString        graphql.String   `graphql:"queryString"`
	QueryStart         graphql.String   `graphql:"queryStart"`
	ThrottleField      graphql.String   `graphql:"throttleField"`
	TimeOfLastTrigger  Long             `graphql:"timeOfLastTrigger"`
	IsStarred          graphql.Boolean  `graphql:"isStarred"`
	Description        graphql.String   `graphql:"description"`
	ThrottleTimeMillis Long             `graphql:"throttleTimeMillis"`
	Enabled            graphql.Boolean  `graphql:"enabled"`
	Actions            []graphql.String `graphql:"actions"`
	Labels             []graphql.String `graphql:"labels"`
	LastError          graphql.String   `graphql:"lastError"`
	QueryOwnership     QueryOwnership   `graphql:"queryOwnership"`
	RunAsUser          struct {
		ID graphql.String `graphql:"id"`
	} `graphql:"runAsUser"`
}

type CreateAlert struct {
	ViewName           graphql.String     `json:"viewName"`
	Name               graphql.String     `json:"name"`
	Description        graphql.String     `json:"description,omitempty"`
	QueryString        graphql.String     `json:"queryString"`
	QueryStart         graphql.String     `json:"queryStart"`
	ThrottleTimeMillis Long               `json:"throttleTimeMillis"`
	ThrottleField      graphql.String     `json:"throttleField,omitempty"`
	RunAsUserID        graphql.String     `json:"runAsUserId,omitempty"`
	Enabled            graphql.Boolean    `json:"enabled"`
	Actions            []graphql.String   `json:"actions"`
	Labels             []graphql.String   `json:"labels"`
	QueryOwnershipType QueryOwnershipType `json:"queryOwnershipType,omitempty"`
}

type UpdateAlert struct {
	ViewName           graphql.String     `json:"viewName"`
	ID                 graphql.String     `json:"id"`
	Name               graphql.String     `json:"name"`
	Description        graphql.String     `json:"description,omitempty"`
	QueryString        graphql.String     `json:"queryString"`
	QueryStart         graphql.String     `json:"queryStart"`
	ThrottleTimeMillis Long               `json:"throttleTimeMillis"`
	ThrottleField      graphql.String     `json:"throttleField,omitempty"`
	RunAsUserID        graphql.String     `json:"runAsUserId,omitempty"`
	Enabled            graphql.Boolean    `json:"enabled"`
	Actions            []graphql.String   `json:"actions"`
	Labels             []graphql.String   `json:"labels"`
	QueryOwnershipType QueryOwnershipType `json:"queryOwnershipType,omitempty"`
}
