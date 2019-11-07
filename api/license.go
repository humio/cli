package api

import (
	"github.com/shurcooL/graphql"
)

type Licenses struct {
	client *Client
}

type License interface {
	ExpiresAt() string
	IssuedAt() string
	LicenseType() string
}

type TrialLicense struct {
	ExpiresAtVal string
	IssuedAtVal  string
}

func (l TrialLicense) LicenseType() string {
	return "trial"
}

func (l TrialLicense) IssuedAt() string {
	return l.IssuedAtVal
}

func (l TrialLicense) ExpiresAt() string {
	return l.ExpiresAtVal
}

type OnPremLicense struct {
	ID            string
	ExpiresAtVal  string
	IssuedAtVal   string
	IssuedTo      string
	NumberOfSeats int
}

func (l OnPremLicense) IssuedAt() string {
	return l.IssuedAtVal
}

func (l OnPremLicense) ExpiresAt() string {
	return l.ExpiresAtVal
}

func (l OnPremLicense) LicenseType() string {
	return "onprem"
}

func (c *Client) Licenses() *Licenses { return &Licenses{client: c} }

func (p *Licenses) Install(license string) error {

	var mutation struct {
		UpdateLicenseKey struct {
			Type string `graphql:"__typename"`
		} `graphql:"updateLicenseKey(license: $license)"`
	}
	variables := map[string]interface{}{
		"license": graphql.String(license),
	}

	return p.client.Mutate(&mutation, variables)
}

func (c *Licenses) Get() (License, error) {
	var query struct {
		License struct {
			ExpiresAt string
			IssuedAt  string
			OnPrem    struct {
				ID       string `graphql:"uid"`
				Owner    string
				MaxUsers int
			} `graphql:"... on OnPremLicense"`
		}
	}

	err := c.client.Query(&query, nil)

	if err != nil {
		return nil, err
	}

	var license License
	if query.License.OnPrem.ID == "" {
		license = TrialLicense{
			ExpiresAtVal: query.License.ExpiresAt,
			IssuedAtVal:  query.License.IssuedAt,
		}
	} else {
		license = OnPremLicense{
			ID:            query.License.OnPrem.ID,
			ExpiresAtVal:  query.License.ExpiresAt,
			IssuedAtVal:   query.License.IssuedAt,
			IssuedTo:      query.License.OnPrem.Owner,
			NumberOfSeats: query.License.OnPrem.MaxUsers,
		}
	}

	return license, nil
}
