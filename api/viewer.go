package api

type Viewer struct {
	client *Client
}

func (c *Client) Viewer() *Viewer { return &Viewer{client: c} }

// Username fetches the username associated with the API Token in use.
func (c *Viewer) Username() (string, error) {
	var query struct {
		Viewer struct {
			Username string
		}
	}

	variables := map[string]interface{}{}

	graphqlErr := c.client.Query(&query, variables)
	return query.Viewer.Username, graphqlErr
}
