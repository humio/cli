package api

import graphql "github.com/cli/shurcooL-graphql"

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Licenses struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type License interface {
	ExpiresAt() string
	IssuedAt() string
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type OnPremLicense struct {
	ID            string
	ExpiresAtVal  string
	IssuedAtVal   string
	IssuedTo      string
	NumberOfSeats int
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (l OnPremLicense) IssuedAt() string {
	return l.IssuedAtVal
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (l OnPremLicense) ExpiresAt() string {
	return l.ExpiresAtVal
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Licenses() *Licenses { return &Licenses{client: c} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (l *Licenses) Install(license string) error {

	var mutation struct {
		UpdateLicenseKey struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"updateLicenseKey(license: $license)"`
	}
	variables := map[string]interface{}{
		"license": graphql.String(license),
	}

	return l.client.Mutate(&mutation, variables)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (l *Licenses) Get() (License, error) {
	var query struct {
		InstalledLicense struct {
			ExpiresAt string
			IssuedAt  string
			OnPrem    struct {
				ID       string `graphql:"uid"`
				Owner    string
				MaxUsers int
			} `graphql:"... on OnPremLicense"`
		}
	}

	err := l.client.Query(&query, nil)
	if err != nil {
		return nil, err
	}

	return OnPremLicense{
		ID:            query.InstalledLicense.OnPrem.ID,
		ExpiresAtVal:  query.InstalledLicense.ExpiresAt,
		IssuedAtVal:   query.InstalledLicense.IssuedAt,
		IssuedTo:      query.InstalledLicense.OnPrem.Owner,
		NumberOfSeats: query.InstalledLicense.OnPrem.MaxUsers,
	}, nil
}
