package api

import (
	graphql "github.com/cli/shurcooL-graphql"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type IngestTokens struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type IngestToken struct {
	Name           string `json:"name"`
	Token          string `json:"token"`
	AssignedParser string `json:"parser"`
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) IngestTokens() *IngestTokens { return &IngestTokens{client: c} }

type ingestTokenData struct {
	Name   string
	Token  string
	Parser *struct {
		Name string
	}
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (i *IngestTokens) List(repo string) ([]IngestToken, error) {
	var query struct {
		Result struct {
			IngestTokens []ingestTokenData
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(repo),
	}

	err := i.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	tokens := make([]IngestToken, len(query.Result.IngestTokens))
	for idx, tokenData := range query.Result.IngestTokens {
		tokens[idx] = *toIngestToken(tokenData)
	}

	return tokens, nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (i *IngestTokens) Get(repoName, tokenName string) (*IngestToken, error) {
	tokensInRepo, err := i.List(repoName)
	if err != nil {
		return nil, err
	}

	for _, token := range tokensInRepo {
		if token.Name == tokenName {
			return &token, nil
		}
	}

	return nil, IngestTokenNotFound(tokenName)
}

func toIngestToken(data ingestTokenData) *IngestToken {
	var parser string
	if data.Parser != nil {
		parser = data.Parser.Name
	}

	return &IngestToken{
		Name:           data.Name,
		Token:          data.Token,
		AssignedParser: parser,
	}
}

func (i *IngestTokens) Add(repositoryName string, tokenName string, parser string, customToken string) (*IngestToken, error) {
	variables := map[string]interface{}{
		"tokenName":      graphql.String(tokenName),
		"repositoryName": graphql.String(repositoryName),
		"parser":         (*graphql.String)(nil),
		"customToken":    (*graphql.String)(nil),
	}

	if parser != "" {
		variables["parser"] = graphql.String(parser)
	}

	if customToken != "" {
		variables["customToken"] = graphql.String(customToken)
	}

	var mutation struct {
		IngestToken struct {
			Name   string
			Token  string
			Parser struct {
				Name string
			}
		} `graphql:"addIngestTokenV3(input: { repositoryName: $repositoryName, name: $tokenName, parser: $parser, customToken: $customToken})"`
	}

	err := i.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	ingestToken := IngestToken{
		Name:           mutation.IngestToken.Name,
		Token:          mutation.IngestToken.Token,
		AssignedParser: mutation.IngestToken.Parser.Name,
	}

	return &ingestToken, nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (i *IngestTokens) Update(repositoryName string, tokenName string, parser string) (*IngestToken, error) {
	if parser == "" {
		var mutation struct {
			Result struct {
				// We have to make a selection, so just take __typename
				Typename graphql.String `graphql:"__typename"`
			} `graphql:"unassignIngestToken(repositoryName: $repositoryName, tokenName: $tokenName)"`
		}

		variables := map[string]interface{}{
			"tokenName":      graphql.String(tokenName),
			"repositoryName": graphql.String(repositoryName),
		}

		err := i.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}
	} else {
		var mutation struct {
			Result struct {
				// We have to make a selection, so just take __typename
				Typename graphql.String `graphql:"__typename"`
			} `graphql:"assignParserToIngestTokenV2(input: { repositoryName: $repositoryName, tokenName: $tokenName, parser: $parser })"`
		}

		variables := map[string]interface{}{
			"tokenName":      graphql.String(tokenName),
			"repositoryName": graphql.String(repositoryName),
			"parser":         graphql.String(parser),
		}

		err2 := i.client.Mutate(&mutation, variables)
		if err2 != nil {
			return nil, err2
		}
	}

	return i.Get(repositoryName, tokenName)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (i *IngestTokens) Remove(repositoryName string, tokenName string) error {
	var mutation struct {
		Result struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"removeIngestToken(repositoryName: $repositoryName, name: $tokenName)"`
	}

	variables := map[string]interface{}{
		"tokenName":      graphql.String(tokenName),
		"repositoryName": graphql.String(repositoryName),
	}

	return i.client.Mutate(&mutation, variables)
}
