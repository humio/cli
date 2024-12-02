package api



type PersonalTokens struct {
	client *Client
}

func (c *Client) PersonalTokens() *PersonalTokens { return &PersonalTokens{client: c} }

func (i *PersonalTokens) Create() (string, error) {
	variables := map[string]interface{}{
	}

	var mutation struct {
		ApiToken string `graphql:"createPersonalUserToken(input: {})"`
	}

	err := i.client.Mutate(&mutation, variables)
	if err != nil {
		return "", err
	}

	return mutation.ApiToken, nil
}
