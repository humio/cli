package api

import (
	"fmt"
	"github.com/SaaldjorMike/graphql"
)

type Roles struct {
	client *Client
}

type Role struct {
	ID                string   `graphql:"id"`
	DisplayName       string   `graphql:"displayName"`
	Color             string   `graphql:"color"`
	Description       string   `graphql:"description`
	ViewPermissions   []string `graphql:"viewPermissions"`
	SystemPermissions []string `graphql:"systemPermissions`
	OrgPermissions    []string `graphql:"organizationPermissions`
}

func (c *Client) Roles() *Roles { return &Roles{client: c} }

func (r *Roles) List() ([]Role, error) {
	var query struct {
		Roles struct {
			Roles []Role
		} `graphql:"roles()"`
	}

	err := r.client.Query(&query, nil)

	var RolesList []Role
	if err == nil {
		RolesList = query.Roles.Roles
	}

	return RolesList, nil
}

func (r *Roles) Create(role *Role) error {
	var mutation struct {
		Role `graphql:"createRole(input: {displayName: $displayName, viewPermissions: $permissions, color: $color, systemPermissions: $systemPermissions, organizationPermissions: $orgPermissions})"`
	}

	viewPermissions := make([]graphql.String, len(role.ViewPermissions))
	for i, permission := range role.ViewPermissions {
		viewPermissions[i] = graphql.String(permission)
	}

	systemPermissions := make([]graphql.String, len(role.SystemPermissions))
	for i, permission := range role.SystemPermissions {
		systemPermissions[i] = graphql.String(permission)
	}

	orgPermissions := make([]graphql.String, len(role.OrgPermissions))
	for i, permission := range role.OrgPermissions {
		orgPermissions[i] = graphql.String(permission)
	}

	variables := map[string]interface{}{
		"displayName":       graphql.String(role.DisplayName),
		"color":             graphql.String(role.Color),
		"description":       graphql.String(role.Description),
		"viewPermissions":   viewPermissions,
		"systemPermissions": systemPermissions,
		"orgPermissions":    orgPermissions,
	}

	return r.client.Mutate(mutation, variables)
}

func (r *Roles) Update(rolename string, newRole *Role) error {
	roleId, err := r.GetRoleID(rolename)
	if roleId == "" || err != nil {
		return fmt.Errorf("unable to find role")
	}

	if newRole == nil {
		return fmt.Errorf("new role values must not be nil")
	}

	var mutation struct {
		Role `graphql:"updateRole(input: {roleId: $roleId, displayName: $displayName, color: $color, description: $description, viewPermissions: $viewPermissions, systemPermissions: $systemPermissions, organizationPermissions: $orgPermissions})"`
	}

	viewPermissions := make([]graphql.String, len(newRole.ViewPermissions))
	for i, permission := range newRole.ViewPermissions {
		viewPermissions[i] = graphql.String(permission)
	}

	systemPermissions := make([]graphql.String, len(newRole.SystemPermissions))
	for i, permission := range newRole.SystemPermissions {
		systemPermissions[i] = graphql.String(permission)
	}

	orgPermissions := make([]graphql.String, len(newRole.OrgPermissions))
	for i, permission := range newRole.OrgPermissions {
		orgPermissions[i] = graphql.String(permission)
	}

	variables := map[string]interface{}{
		"roleId":                  graphql.String(roleId),
		"displayName":             graphql.String(newRole.DisplayName),
		"color":                   graphql.String(newRole.Color),
		"description":             graphql.String(newRole.Description),
		"viewPermissions":         viewPermissions,
		"systemPermissions":       systemPermissions,
		"organizationPermissions": orgPermissions,
	}

	return r.client.Mutate(mutation, variables)
}

func (r *Roles) RemoveRole(rolename string) error {
	var mutation struct {
		RemoveRole struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"removeRole(input: {roleId: $roleId})"`
	}

	role, err := r.client.Roles().Get(rolename)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"roleId": graphql.String(role.ID),
	}

	return r.client.Mutate(mutation, variables)
}

func (r *Roles) Get(rolename string) (*Role, error) {
	roleId, err := r.GetRoleID(rolename)
	if roleId == "" || err != nil {
		return nil, fmt.Errorf("unable to get role id")
	}

	var query struct {
		Role `graphql:"role(roleId: $roleId)"`
	}

	variables := map[string]interface{}{
		"roleId": graphql.String(roleId),
	}

	err = r.client.Query(query, variables)
	if err != nil {
		return nil, err
	}

	return &query.Role, nil
}

func (r *Roles) GetRoleID(rolename string) (string, error) {
	roles, err := r.List()
	if err != nil {
		return "", fmt.Errorf("unable to list roles: %w", err)
	}
	var roleId string
	for _, role := range roles {
		if role.DisplayName == rolename {
			roleId = role.ID
			break
		}
	}

	return roleId, nil
}
