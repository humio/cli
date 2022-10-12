package api

import (
	"errors"

	"github.com/shurcooL/graphql"
)

type Groups struct {
	client *Client
}

type Group struct {
	ID          string
	DisplayName string
}

type GroupChangeSet struct {
	DefaultQueryPrefix *string
	DefaultRole        *Role
	// searchDomainRoles                []SearchDomainRole
	// reconcileSearchDomainQueryPrefix bool
	GroupID *string
}

func (c *Client) Groups() *Groups { return &Groups{client: c} }

var ErrUserNotFound = errors.New("user not found")

func (g *Groups) Find(groupName string) (*Group, error) {
	var query struct {
		Page struct {
			Groups []Group `graphql:"page"`
		} `graphql:"groupsPage(search: $groupName, pageNumber:1, pageSize:1)"`
	}

	variables := map[string]interface{}{
		"groupName": graphql.String(groupName),
	}

	err := g.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	return &query.Page.Groups[0], nil
}

func (g *Groups) List() ([]Group, error) {
	var query struct {
		Page struct {
			Groups []Group `graphql:"page"`
		} `graphql:"groupsPage(pageNumber:1,pageSize:2147483647)"`
	}

	err := g.client.Query(&query, nil)
	if err != nil {
		return nil, err
	}

	return query.Page.Groups, nil
}

// TODO:
// A way to link a list of repos/views to a group, which also contains a link to a corresponding role.
// A way to link a list of users to a group.
// I suspect the way we would like to "link" these things is by referencing the names of the resources; however, it's possible it may not be that simple, or that there may be a built-in way to do this. For example, how the views manage connections to repos https://github.com/humio/humio-operator/blob/master/api/v1alpha1/humioview_types.go#L53-L54. There may be cases where we have to pass a generated ID as the link, but then first have to look up the resource by name to fetch the ID to make it human-usable.

func (g *Groups) Create(groupName string) error {
	var mutation struct {
		addGroup struct {
			Group struct {
				ID string
			}
		} `graphql:"addGroup(displayName: $groupName)"`
	}

	variables := map[string]interface{}{
		"groupName": graphql.String(groupName),
	}

	err := g.client.Mutate(&mutation, variables)
	if err != nil {
		return err
	}

	return nil
}

func (g *Groups) Delete(groupID string) error {
	var mutation struct {
		removeGroup struct {
			Group struct {
				ID string
			}
		} `graphql:"removeGroup(groupId: $groupID)"`
	}

	variables := map[string]interface{}{
		"groupID": graphql.String(groupID),
	}

	err := g.client.Mutate(&mutation, variables)
	if err != nil {
		return err
	}

	return nil
}

func (g *Groups) Update(displayName string, changeset GroupChangeSet) error {
	var mutation struct {
		updateGroup struct {
			Group struct {
				ID string
			}
		} `graphql:"updateGroup(input:{groupId:$groupID, displayName: $displayName})"`
	}

	variables := map[string]interface{}{
		"groupID":     graphql.String(*changeset.GroupID),
		"displayName": graphql.String(displayName),
	}

	err := g.client.Mutate(&mutation, variables)
	if err != nil {
		return err
	}

	return nil
}

func (g *Groups) AddUserToGroup(groupID string, userID string) error {
	var mutation struct {
		AddUsersToGroup struct {
			Group struct {
				Users []struct {
					ID string
				}
			}
		} `graphql:"addUsersToGroup(input:{users:[$userID], groupId: $groupID})"`
	}

	variables := map[string]interface{}{
		"userID":  graphql.String(userID),
		"groupID": graphql.String(groupID),
	}

	err := g.client.Mutate(&mutation, variables)
	if err != nil {
		return err
	}

	var found bool
	for _, user := range mutation.AddUsersToGroup.Group.Users {
		if user.ID == userID {
			found = true
			break
		}
	}

	if !found {
		return ErrUserNotFound
	}

	return nil
}

func (g *Groups) RemoveUserFromGroup(groupID string, userID string) error {
	var mutation struct {
		RemoveUsersFromGroup struct {
			Group struct {
				ID string
			}
		} `graphql:"removeUsersFromGroup(input:{users:[$userID], groupId: $groupID})"`
	}

	variables := map[string]interface{}{
		"userID":  graphql.String(userID),
		"groupID": graphql.String(groupID),
	}

	return g.client.Mutate(&mutation, variables)
}
