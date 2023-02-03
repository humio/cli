package api

import (
	"github.com/SaaldjorMike/graphql"
)

type Licenses struct {
	client *Client
}

type License interface {
	ExpiresAt() string
	IssuedAt() string
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

func (c *Client) Licenses() *Licenses { return &Licenses{client: c} }

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
