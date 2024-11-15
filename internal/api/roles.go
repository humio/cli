package api

import (
	"context"
	"fmt"

	"github.com/humio/cli/internal/api/humiographql"
)

type Roles struct {
	client *Client
}

type Role struct {
	ID                string
	DisplayName       string
	ViewPermissions   []string
	SystemPermissions []string
	OrgPermissions    []string
}

func (c *Client) Roles() *Roles { return &Roles{client: c} }

func (r *Roles) List() ([]Role, error) {
	resp, err := humiographql.ListRoles(context.Background(), r.client)
	if err != nil {
		return nil, err
	}

	respRoles := resp.GetRoles()
	roles := make([]Role, len(respRoles))
	for idx, role := range respRoles {
		respViewPermissions := role.GetViewPermissions()
		viewPermissions := make([]string, len(respViewPermissions))
		for k, perm := range respViewPermissions {
			viewPermissions[k] = string(perm)
		}

		respOrgPermissions := role.GetOrganizationPermissions()
		orgPermissions := make([]string, len(respOrgPermissions))
		for k, perm := range respOrgPermissions {
			orgPermissions[k] = string(perm)
		}

		respSystemPermissions := role.GetSystemPermissions()
		systemPermissions := make([]string, len(respSystemPermissions))
		for k, perm := range respSystemPermissions {
			systemPermissions[k] = string(perm)
		}

		roles[idx] = Role{
			ID:                role.GetId(),
			DisplayName:       role.GetDisplayName(),
			ViewPermissions:   viewPermissions,
			OrgPermissions:    orgPermissions,
			SystemPermissions: systemPermissions,
		}
	}

	return roles, nil
}

func (r *Roles) Create(role *Role) error {
	viewPermissions := make([]humiographql.Permission, len(role.ViewPermissions))
	for k, perm := range role.ViewPermissions {
		viewPermissions[k] = humiographql.Permission(perm)
	}
	orgPermissions := make([]humiographql.OrganizationPermission, len(role.OrgPermissions))
	for k, perm := range role.OrgPermissions {
		orgPermissions[k] = humiographql.OrganizationPermission(perm)
	}
	systemPermissions := make([]humiographql.SystemPermission, len(role.SystemPermissions))
	for k, perm := range role.SystemPermissions {
		systemPermissions[k] = humiographql.SystemPermission(perm)
	}

	_, err := humiographql.CreateRole(context.Background(), r.client, role.DisplayName, viewPermissions, orgPermissions, systemPermissions)
	return err
}

func (r *Roles) Update(rolename string, newRole *Role) error {
	roleId, err := r.GetRoleID(rolename)
	if roleId == "" || err != nil {
		return fmt.Errorf("unable to find role")
	}

	if newRole == nil {
		return fmt.Errorf("new role values must not be nil")
	}

	viewPermissions := make([]humiographql.Permission, len(newRole.ViewPermissions))
	for k, perm := range newRole.ViewPermissions {
		viewPermissions[k] = humiographql.Permission(perm)
	}
	orgPermissions := make([]humiographql.OrganizationPermission, len(newRole.OrgPermissions))
	for k, perm := range newRole.OrgPermissions {
		orgPermissions[k] = humiographql.OrganizationPermission(perm)
	}
	systemPermissions := make([]humiographql.SystemPermission, len(newRole.SystemPermissions))
	for k, perm := range newRole.SystemPermissions {
		systemPermissions[k] = humiographql.SystemPermission(perm)
	}

	_, err = humiographql.UpdateRole(context.Background(), r.client, roleId, newRole.DisplayName, viewPermissions, orgPermissions, systemPermissions)
	return err
}

func (r *Roles) RemoveRole(rolename string) error {
	role, err := r.client.Roles().Get(rolename)
	if err != nil {
		return err
	}

	_, err = humiographql.RemoveRoleByID(context.Background(), r.client, role.ID)
	return err
}

func (r *Roles) Get(rolename string) (*Role, error) {
	roleId, err := r.GetRoleID(rolename)
	if roleId == "" || err != nil {
		return nil, fmt.Errorf("unable to get role id")
	}

	resp, err := humiographql.GetRoleByID(context.Background(), r.client, roleId)
	if err != nil {
		return nil, err
	}
	role := resp.GetRole()
	respViewPermissions := role.GetViewPermissions()
	viewPermissions := make([]string, len(respViewPermissions))
	for k, perm := range respViewPermissions {
		viewPermissions[k] = string(perm)
	}

	respOrgPermissions := role.GetOrganizationPermissions()
	orgPermissions := make([]string, len(respOrgPermissions))
	for k, perm := range respOrgPermissions {
		orgPermissions[k] = string(perm)
	}

	respSystemPermissions := role.GetSystemPermissions()
	systemPermissions := make([]string, len(respSystemPermissions))
	for k, perm := range respSystemPermissions {
		systemPermissions[k] = string(perm)
	}

	return &Role{
		ID:                role.GetId(),
		DisplayName:       role.GetDisplayName(),
		ViewPermissions:   viewPermissions,
		OrgPermissions:    orgPermissions,
		SystemPermissions: systemPermissions,
	}, nil
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
