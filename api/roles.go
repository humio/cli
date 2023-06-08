package api

import (
	"fmt"
)

type Roles struct {
	client *Client
}

type Role struct {
	ID                string   `graphql:"id"`
	DisplayName       string   `graphql:"displayName"`
	Color             string   `graphql:"color"`
	Description       string   `graphql:"description"`
	ViewPermissions   []string `graphql:"viewPermissions"`
	SystemPermissions []string `graphql:"systemPermissions"`
	OrgPermissions    []string `graphql:"organizationPermissions"`
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

	viewPermissions := make([]string, len(role.ViewPermissions))
	copy(viewPermissions, role.ViewPermissions)

	systemPermissions := make([]string, len(role.SystemPermissions))
	copy(systemPermissions, role.SystemPermissions)

	orgPermissions := make([]string, len(role.OrgPermissions))
	copy(orgPermissions, role.OrgPermissions)

	variables := map[string]interface{}{
		"displayName":       role.DisplayName,
		"color":             role.Color,
		"description":       role.Description,
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

	viewPermissions := make([]string, len(newRole.ViewPermissions))
	copy(viewPermissions, newRole.ViewPermissions)

	systemPermissions := make([]string, len(newRole.SystemPermissions))
	copy(systemPermissions, newRole.SystemPermissions)

	orgPermissions := make([]string, len(newRole.OrgPermissions))
	copy(orgPermissions, newRole.OrgPermissions)

	variables := map[string]interface{}{
		"roleId":                  roleId,
		"displayName":             newRole.DisplayName,
		"color":                   newRole.Color,
		"description":             newRole.Description,
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
			Typename string `graphql:"__typename"`
		} `graphql:"removeRole(input: {roleId: $roleId})"`
	}

	role, err := r.client.Roles().Get(rolename)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"roleId": role.ID,
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
		"roleId": roleId,
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
