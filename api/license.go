package api

import (
	"github.com/shurcooL/graphql"
)

type License struct {
	client *Client
}

type LicenseData struct {
	ExpiresAt string
	IssuedAt  string
}

func (c *Client) License() *License { return &License{client: c} }

func (p *License) Install(license string) error {

	var mutation struct {
		CreateParser struct {
			Type string `graphql:"__typename"`
		} `graphql:"updateLicenseKey(license: $license)"`
	}
	variables := map[string]interface{}{
		"license": graphql.String(license),
	}

	return p.client.Mutate(&mutation, variables)
}

func (c *License) Get() (LicenseData, error) {
	var query struct {
		License LicenseData
	}
	variables := map[string]interface{}{}

	err := c.client.Query(&query, variables)

	return query.License, err
}
