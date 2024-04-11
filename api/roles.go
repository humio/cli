package api

import (
	"fmt"

	graphql "github.com/cli/shurcooL-graphql"
)

type Roles struct {
	client *Client
}

func (c *Client) Roles() *Roles { return &Roles{client: c} }

type Role struct {
	ID                string   `json:"id"`
	DisplayName       string   `json:"displayName"`
	Color             string   `json:"color,omitempty"`
	Description       string   `json:"description,omitempty"`
	ViewPermissions   []string `json:"viewPermissions"`
	SystemPermissions []string `json:"systemPermissions,omitempty"`
	OrgPermissions    []string `graphql:"organizationPermissions" json:"organizationPermissions,omitempty"`
	GroupsCount       int      `json:"groupsCount"`
	UsersCount        int      `json:"usersCount"`
}

// List returns a list of roles in the Humio instance.
func (r *Roles) List() ([]Role, error) {
	var query struct {
		Roles []Role `graphql:"roles"`
	}

	err := r.client.Query(&query, nil)

	if err != nil {
		return nil, err
	}

	return query.Roles, nil
}

// Create adds a new role to the Humio instance.
func (r *Roles) Create(role AddRoleInput) error {
	var mutation struct {
		Results struct {
			Role struct {
				ID string
			}
		} `graphql:"createRole(input:$input)"`
	}
	variables := map[string]interface{}{
		"input": role,
	}
	return r.client.Mutate(&mutation, variables)
}

// Upddate updates a role in the Humio instance.
func (r *Roles) Update(role UpdateRoleInput) error {
	var mutation struct {
		Results struct {
			Role struct {
				ID string
			}
		} `graphql:"updateRole(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": role,
	}

	return r.client.Mutate(&mutation, variables)
}

// Delete removes a role from the Humio instance.
func (r *Roles) Delete(name string) error {
	var mutation struct {
		Result struct {
			Result bool
		} `graphql:"removeRole(roleId: $roleId)"`
	}

	role, err := r.client.Roles().Get(name)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"roleId": graphql.String(role.ID),
	}
	return r.client.Mutate(&mutation, variables)
}

// Get returns a role given its name.
func (r *Roles) Get(name string) (*Role, error) {
	roleId, err := r.GetRoleID(name)
	if roleId == "" || err != nil {
		return nil, fmt.Errorf("unable to get role id")
	}

	var query struct {
		Role Role `graphql:"role(roleId: $roleId)"`
	}

	variables := map[string]interface{}{
		"roleId": graphql.String(roleId),
	}

	err = r.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}
	return &query.Role, nil
}

// GetRoleID returns the ID of a role given its name.
func (r *Roles) GetRoleID(name string) (string, error) {
	roles, err := r.List()
	if err != nil {
		return "", fmt.Errorf("unable to list roles: %w", err)
	}
	var roleId string
	for _, role := range roles {
		if role.DisplayName == name {
			roleId = string(role.ID)
			break
		}
	}
	return roleId, nil
}
