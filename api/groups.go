package api

import (
	"errors"
)

type Groups struct {
	client *Client
}

type Group struct {
	ID          string
	DisplayName string
}

func (c *Client) Groups() *Groups { return &Groups{client: c} }

var ErrUserNotFound = errors.New("user not found")

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
		"userID":  userID,
		"groupID": groupID,
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
		"userID":  userID,
		"groupID": groupID,
	}

	return g.client.Mutate(&mutation, variables)
}
