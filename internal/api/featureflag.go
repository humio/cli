package api

import (
	"context"
	"sort"

	"github.com/humio/cli/internal/api/humiographql"
)

type FeatureFlagName string

type FeatureFlag struct {
	Flag         FeatureFlagName
	Experimental bool
	Description  string
}

type FeatureFlags struct {
	client *Client
}

func (c *Client) FeatureFlags() *FeatureFlags {
	return &FeatureFlags{client: c}
}

func (f *FeatureFlags) SupportedFlags() ([]FeatureFlag, error) {
	resp, err := humiographql.GetSupportedFeatureFlags(context.Background(), f.client)
	if err != nil {
		return nil, err
	}

	respFeatureFlags := resp.GetFeatureFlags()
	supportedFlags := make([]FeatureFlag, len(respFeatureFlags))
	for idx, flag := range respFeatureFlags {
		supportedFlags[idx] = FeatureFlag{
			Flag:         FeatureFlagName(flag.GetFlag()),
			Experimental: flag.GetExperimental(),
			Description:  flag.GetDescription(),
		}
	}

	sort.Slice(supportedFlags, func(i, j int) bool {
		return supportedFlags[i].Flag < supportedFlags[j].Flag
	})
	return supportedFlags, nil
}

func (f *FeatureFlags) EnableGlobally(flag FeatureFlagName) error {
	_, err := humiographql.EnableFeatureFlagGlobally(context.Background(), f.client, humiographql.FeatureFlag(flag))
	return err
}

func (f *FeatureFlags) DisableGlobally(flag FeatureFlagName) error {
	_, err := humiographql.DisableFeatureFlagGlobally(context.Background(), f.client, humiographql.FeatureFlag(flag))
	return err
}

func (f *FeatureFlags) EnableForOrganization(organizationID string, flag FeatureFlagName) error {
	_, err := humiographql.EnableFeatureFlagForOrganization(context.Background(), f.client, humiographql.FeatureFlag(flag), organizationID)
	return err
}

func (f *FeatureFlags) DisableForOrganization(organizationID string, flag FeatureFlagName) error {
	_, err := humiographql.DisableFeatureFlagForOrganization(context.Background(), f.client, humiographql.FeatureFlag(flag), organizationID)
	return err
}

func (f *FeatureFlags) EnableForUser(userID string, flag FeatureFlagName) error {
	_, err := humiographql.EnableFeatureFlagForUser(context.Background(), f.client, humiographql.FeatureFlag(flag), userID)
	return err
}

func (f *FeatureFlags) DisableForUser(userID string, flag FeatureFlagName) error {
	_, err := humiographql.DisableFeatureFlagForUser(context.Background(), f.client, humiographql.FeatureFlag(flag), userID)
	return err
}
