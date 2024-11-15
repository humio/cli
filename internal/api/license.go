package api

import (
	"context"
	"fmt"

	"github.com/humio/cli/internal/api/humiographql"
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
	NumberOfSeats *int
}

func (l OnPremLicense) IssuedAt() string {
	return l.IssuedAtVal
}

func (l OnPremLicense) ExpiresAt() string {
	return l.ExpiresAtVal
}

func (c *Client) Licenses() *Licenses { return &Licenses{client: c} }

func (l *Licenses) Install(license string) error {
	_, err := humiographql.UpdateLicenseKey(context.Background(), l.client, license)
	return err
}

func (l *Licenses) Get() (License, error) {
	resp, err := humiographql.GetLicense(context.Background(), l.client)
	if err != nil {
		return nil, err
	}

	installedLicense := resp.GetInstalledLicense()
	if installedLicense == nil {
		return nil, fmt.Errorf("no license installed")
	}
	switch v := (*installedLicense).(type) {
	case *humiographql.GetLicenseInstalledLicenseOnPremLicense:
		return OnPremLicense{
			ExpiresAtVal:  v.GetExpiresAt().String(),
			IssuedAtVal:   v.GetIssuedAt().String(),
			ID:            v.GetUid(),
			IssuedTo:      v.GetOwner(),
			NumberOfSeats: v.GetMaxUsers(),
		}, nil
	}

	return nil, fmt.Errorf("unsupported license type")
}
