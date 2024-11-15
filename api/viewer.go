package api

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Viewer struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Viewer() *Viewer { return &Viewer{client: c} }

// Username fetches the username associated with the API Token in use.
// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Viewer) Username() (string, error) {
	var query struct {
		Viewer struct {
			Username string
		}
	}

	err := c.client.Query(&query, nil)
	return query.Viewer.Username, err
}

// ApiToken fetches the api token for the user who is currently authenticated.
// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Viewer) ApiToken() (string, error) {
	var query struct {
		Viewer struct {
			ApiToken string
		}
	}

	err := c.client.Query(&query, nil)
	return query.Viewer.ApiToken, err
}
