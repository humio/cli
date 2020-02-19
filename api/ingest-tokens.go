package api

import (
	"fmt"

	"github.com/shurcooL/graphql"
)

type IngestTokens struct {
	client *Client
}

type IngestToken struct {
	Name           string `json:"name"`
	Token          string `json:"token"`
	AssignedParser string `json:"parser"`
}

func (c *Client) IngestTokens() *IngestTokens { return &IngestTokens{client: c} }

type ingestTokenData struct {
	Name   string
	Token  string
	Parser *struct {
		Name string
	}
}

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

	return nil, fmt.Errorf("could not find an ingest token with name '%s' in repo '%s'", tokenName, repoName)
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

func (i *IngestTokens) Add(repo string, name string, parserName string) (*IngestToken, error) {
	var mutation struct {
		Result struct {
			IngestToken ingestTokenData
		} `graphql:"addIngestToken(repositoryName: $repositoryName, name: $name, parser: $parser)"`
	}

	var parserNameVar graphql.String
	if parserName != "" {
		parserNameVar = graphql.String(parserName)
	}
	variables := map[string]interface{}{
		"name":           graphql.String(name),
		"repositoryName": graphql.String(repo),
		"parser":         parserNameVar,
	}

	err := i.client.Mutate(&mutation, variables)

	if err != nil {
		return nil, err
	}

	return toIngestToken(mutation.Result.IngestToken), err
}

func (i *IngestTokens) Remove(repo string, tokenName string) error {
	var mutation struct {
		Result struct {
			Type string `graphql:"__typename"`
		} `graphql:"removeIngestToken(repositoryName: $repositoryName, name: $name)"`
	}

	variables := map[string]interface{}{
		"name":           graphql.String(tokenName),
		"repositoryName": graphql.String(repo),
	}

	return i.client.Mutate(&mutation, variables)
}
