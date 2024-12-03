package api

import graphql "github.com/cli/shurcooL-graphql"

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type FeatureFlag string

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type FeatureFlags struct {
	c *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) FeatureFlags() *FeatureFlags {
	return &FeatureFlags{c: c}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (f *FeatureFlags) SupportedFlags() ([]FeatureFlag, error) {
	var query struct {
		Type struct {
			EnumValues []struct {
				Name string
			} `graphql:"enumValues"`
		} `graphql:"__type(name: \"FeatureFlag\")"`
	}

	err := f.c.Query(&query, nil)
	if err != nil {
		return nil, err
	}

	var result []FeatureFlag
	for _, flag := range query.Type.EnumValues {
		result = append(result, FeatureFlag(flag.Name))
	}

	return result, nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (f *FeatureFlags) EnableGlobally(flag FeatureFlag) error {
	var mutation struct {
		EnableFeature bool `graphql:"enableFeature(feature: $feature)"`
	}

	variables := map[string]interface{}{
		"feature": flag,
	}

	return f.c.Mutate(&mutation, variables)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (f *FeatureFlags) DisableGlobally(flag FeatureFlag) error {
	var mutation struct {
		DisableFeature bool `graphql:"disableFeature(feature: $feature)"`
	}

	variables := map[string]interface{}{
		"feature": flag,
	}

	return f.c.Mutate(&mutation, variables)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (f *FeatureFlags) EnableForOrganization(organizationID string, flag FeatureFlag) error {
	var mutation struct {
		EnableFeature bool `graphql:"enableFeatureForOrg(feature: $feature, orgId: $orgId)"`
	}

	variables := map[string]interface{}{
		"feature": flag,
		"orgId":   graphql.String(organizationID),
	}

	return f.c.Mutate(&mutation, variables)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (f *FeatureFlags) DisableForOrganization(organizationID string, flag FeatureFlag) error {
	var mutation struct {
		DisableFeature bool `graphql:"disableFeatureForOrg(feature: $feature, orgId: $orgId)"`
	}

	variables := map[string]interface{}{
		"feature": flag,
		"orgId":   graphql.String(organizationID),
	}

	return f.c.Mutate(&mutation, variables)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (f *FeatureFlags) EnableForUser(userID string, flag FeatureFlag) error {
	var mutation struct {
		EnableFeature bool `graphql:"enableFeatureForUser(feature: $feature, userId: $userId)"`
	}

	variables := map[string]interface{}{
		"feature": flag,
		"userId":  graphql.String(userID),
	}

	return f.c.Mutate(&mutation, variables)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (f *FeatureFlags) DisableForUser(userID string, flag FeatureFlag) error {
	var mutation struct {
		DisableFeature bool `graphql:"disableFeatureForUser(feature: $feature, userId: $userId)"`
	}

	variables := map[string]interface{}{
		"feature": flag,
		"userId":  graphql.String(userID),
	}

	return f.c.Mutate(&mutation, variables)
}
