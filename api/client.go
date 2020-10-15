package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/shurcooL/graphql"
)

type Client struct {
	config Config
}

type Config struct {
	Address       *url.URL
	Token         string
	CACertificate []byte
	Insecure      bool
}

func DefaultConfig() Config {
	return Config{}
}

func (c *Client) Address() *url.URL {
	return c.config.Address
}

func (c *Client) Token() string {
	return c.config.Token
}

func (c *Client) CACertificate() []byte {
	return c.config.CACertificate
}

func (c *Client) Insecure() bool {
	return c.config.Insecure
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) newGraphQLClient() *graphql.Client {
	httpClient := c.newHTTPClientWithHeaders(map[string]string{
		"Authorization": "Bearer " + c.Token(),
	})
	graphqlURL, err := c.Address().Parse("/graphql")
	if err != nil {
		panic(err)
	}
	return graphql.NewClient(graphqlURL.String(), httpClient)
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

// JSONContentType is "application/json"
const JSONContentType string = "application/json"

func (c *Client) HTTPRequest(httpMethod string, path string, body io.Reader) (*http.Response, error) {
	return c.HTTPRequestContext(context.Background(), httpMethod, path, body, JSONContentType)
}

func (c *Client) HTTPRequestContext(ctx context.Context, httpMethod string, path string, body io.Reader, contentType string) (*http.Response, error) {
	if body == nil {
		body = bytes.NewReader(nil)
	}

	url, err := c.Address().Parse(path)
	if err != nil {
		return nil, err
	}

	req, reqErr := http.NewRequestWithContext(ctx, httpMethod, url.String(), body)
	if reqErr != nil {
		return nil, reqErr
	}

	var client = c.newHTTPClientWithHeaders(map[string]string{
		"Authorization": "Bearer " + c.Token(),
		"Content-Type":  contentType,
	})
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
