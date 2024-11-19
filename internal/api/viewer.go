package api

import (
	"context"

	"github.com/humio/cli/internal/api/humiographql"
)

type Viewer struct {
	client *Client
}

func (c *Client) Viewer() *Viewer { return &Viewer{client: c} }

// Username fetches the username associated with the API Token in use.
func (c *Viewer) Username() (string, error) {
	resp, err := humiographql.GetUsername(context.Background(), c.client)
	if err != nil {
		return "", err
	}
	viewer := resp.GetViewer()
	return viewer.GetUsername(), nil
}
