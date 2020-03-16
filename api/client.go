package api

import (
	"bytes"
	"context"
	"net/http"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type Client struct {
	config Config
}

type Config struct {
	Address string
	Token   string
}

func DefaultConfig() Config {
	config := Config{
		Address: "",
		Token:   "",
	}

	return config
}

func (c *Client) Address() string {
	return c.config.Address
}

func (c *Client) Token() string {
	return c.config.Token
}

func NewClient(config Config) (*Client, error) {
	return &Client{
		config: config,
	}, nil
}

func (c *Client) newGraphQLClient() *graphql.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.config.Token},
	)

	httpClient := oauth2.NewClient(context.Background(), src)
	return graphql.NewClient(c.Address()+"graphql", httpClient)
}

func (c *Client) Query(query interface{}, variables map[string]interface{}) error {
	client := c.newGraphQLClient()
	graphqlErr := client.Query(context.Background(), query, variables)
	return graphqlErr
}

func (c *Client) Mutate(mutation interface{}, variables map[string]interface{}) error {
	client := c.newGraphQLClient()
	graphqlErr := client.Mutate(context.Background(), mutation, variables)
	return graphqlErr
}

func (c *Client) HTTPRequest(httpMethod string, path string, body *bytes.Buffer) (*http.Response, error) {
	return c.HTTPRequestContext(context.Background(), httpMethod, path, body)
}

func (c *Client) HTTPRequestContext(ctx context.Context, httpMethod string, path string, body *bytes.Buffer) (*http.Response, error) {
	if body == nil {
		body = bytes.NewBuffer([]byte(""))
	}

	url := c.Address() + path

	req, reqErr := http.NewRequestWithContext(ctx, httpMethod, url, body)
	req.Header.Set("Authorization", "Bearer "+c.Token())
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{}

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}

func optBoolArg(v *bool) *graphql.Boolean {
	var argPtr *graphql.Boolean
	if v != nil {
		argPtr = graphql.NewBoolean(graphql.Boolean(*v))
	}
	return argPtr
}

func optStringArg(v *string) *graphql.String {
	var argPtr *graphql.String
	if v != nil {
		argPtr = graphql.NewString(graphql.String(*v))
	}
	return argPtr
}
