package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	graphql "github.com/cli/shurcooL-graphql"
)

const defaultUserAgent = "Humio-go-client/unknown"

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Client struct {
	config        Config
	httpTransport *http.Transport
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Config struct {
	Address           *url.URL
	UserAgent         string
	Token             string
	CACertificatePEM  string
	Insecure          bool
	ProxyOrganization string
	DialContext       func(ctx context.Context, network, addr string) (net.Conn, error)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func DefaultConfig() Config {
	config := Config{}

	return config
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Address() *url.URL {
	return c.config.Address
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Token() string {
	return c.config.Token
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) CACertificate() string {
	return c.config.CACertificatePEM
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Insecure() bool {
	return c.config.Insecure
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Config() Config {
	return c.config
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func NewClient(config Config) *Client {
	httpTransport := NewHttpTransport(config)
	return NewClientWithTransport(config, httpTransport)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func NewClientWithTransport(config Config, httpTransport *http.Transport) *Client {
	if config.Address != nil && !strings.HasSuffix(config.Address.Path, "/") {
		config.Address.Path = config.Address.Path + "/"
	}

	if config.UserAgent == "" {
		config.UserAgent = defaultUserAgent
	}

	return &Client{
		config:        config,
		httpTransport: httpTransport,
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) headers() map[string]string {
	headers := map[string]string{}

	if c.Token() != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", c.Token())
	}

	if c.config.ProxyOrganization != "" {
		headers["ProxyOrganization"] = c.config.ProxyOrganization
	}

	if c.config.UserAgent != "" {
		headers["User-Agent"] = c.config.UserAgent
	}

	return headers
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) newGraphQLClient() (*graphql.Client, error) {
	httpClient := c.newHTTPClientWithHeaders(c.headers())
	graphqlURL, err := c.Address().Parse("graphql")
	if err != nil {
		return nil, err
	}
	return graphql.NewClient(graphqlURL.String(), httpClient), nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Query(query interface{}, variables map[string]interface{}) error {
	client, err := c.newGraphQLClient()
	if err != nil {
		return err
	}
	return client.Query(context.Background(), query, variables)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Mutate(mutation interface{}, variables map[string]interface{}) error {
	client, err := c.newGraphQLClient()
	if err != nil {
		return err
	}
	return client.Mutate(context.Background(), mutation, variables)
}

// JSONContentType is "application/json"
// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
const JSONContentType string = "application/json"

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
const ZIPContentType string = "application/zip"

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) HTTPRequest(httpMethod string, path string, body io.Reader) (*http.Response, error) {
	return c.HTTPRequestContext(context.Background(), httpMethod, path, body, JSONContentType)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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

	headers := c.headers()
	headers["Content-Type"] = contentType

	var client = c.newHTTPClientWithHeaders(headers)
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
