package api

import (
	"context"

	"github.com/humio/cli/internal/api/humiographql"
)

type Tokens struct {
	client *Client
}

func (c *Client) Tokens() *Tokens { return &Tokens{client: c} }

func (t *Tokens) Rotate(tokenID string) (string, error) {
	resp, err := humiographql.RotateTokenByID(context.Background(), t.client, tokenID)
	if err != nil {
		return "", err
	}

	return resp.GetRotateToken(), nil
}
