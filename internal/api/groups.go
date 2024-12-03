package api

import (
	"context"

	"github.com/humio/cli/internal/api/humiographql"
)

type Groups struct {
	client *Client
}

type Group struct {
	ID          string
	DisplayName string
}

func (c *Client) Groups() *Groups { return &Groups{client: c} }

func (g *Groups) List() ([]Group, error) {
	resp, err := humiographql.ListGroups(context.Background(), g.client)
	if err != nil {
		return nil, err
	}
	respGroups := resp.GetGroupsPage()
	respGroupsPage := respGroups.GetPage()
	groups := make([]Group, len(respGroupsPage))
	for idx, group := range respGroupsPage {
		groups[idx] = Group{
			ID:          group.GetId(),
			DisplayName: group.GetDisplayName(),
		}
	}

	return groups, nil
}

func (g *Groups) AddUserToGroup(groupID string, userID string) error {
	_, err := humiographql.AddUserToGroup(context.Background(), g.client, groupID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (g *Groups) RemoveUserFromGroup(groupID string, userID string) error {
	_, err := humiographql.RemoveUserFromGroup(context.Background(), g.client, groupID, userID)
	if err != nil {
		return err
	}

	return nil
}
