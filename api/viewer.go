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

	graphqlErr := c.client.Query(&query, nil)
	return query.Viewer.Username, graphqlErr
}

// ApiToken fetches the api token for the user who is currently authenticated.
func (c *Viewer) ApiToken() (string, error) {
	var query struct {
		Viewer struct {
			ApiToken string
		}
	}

	graphqlErr := c.client.Query(&query, nil)
	return query.Viewer.ApiToken, graphqlErr
}
