package humiographql

import graphql "github.com/cli/shurcooL-graphql"

type QueryOwnershipTypeName string

const (
	QueryOwnershipTypeNameOrganization QueryOwnershipTypeName = "OrganizationOwnership"
	QueryOwnershipTypeNameUser         QueryOwnershipTypeName = "UserOwnership"
)

type QueryOwnershipType string

const (
	QueryOwnershipTypeUser         QueryOwnershipType = "User"
	QueryOwnershipTypeOrganization QueryOwnershipType = "Organization"
)

type QueryOwnership struct {
	ID                     graphql.String         `graphql:"id"`
	QueryOwnershipTypeName QueryOwnershipTypeName `graphql:"__typename"`
}
