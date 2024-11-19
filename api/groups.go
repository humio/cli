package api

import (
	"errors"

	graphql "github.com/cli/shurcooL-graphql"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Groups struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Group struct {
	ID          string
	DisplayName string
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Groups() *Groups { return &Groups{client: c} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
var ErrUserNotFound = errors.New("user not found")

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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
