package api

import (
	"context"
	"os"

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

	if addr := os.Getenv("HUMIO_ADDR"); addr != "" {
		config.Address = addr
	}

	if token := os.Getenv("HUMIO_API_TOKEN"); token != "" {
		config.Token = token
	}

	return config
}

func (c *Client) Address() string {
	return c.config.Address
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
