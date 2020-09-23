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

func (i *IngestTokens) Add(repositoryName string, tokenName string, parserName string) (*IngestToken, error) {
	var mutation struct {
		Result struct {
			IngestToken ingestTokenData
		} `graphql:"addIngestToken(repositoryName: $repositoryName, name: $tokenName, parser: $parserName)"`
	}

	variables := map[string]interface{}{
		"tokenName":      graphql.String(tokenName),
		"repositoryName": graphql.String(repositoryName),
		"parserName":     graphql.String(parserName),
	}

	err := i.client.Mutate(&mutation, variables)

	if err != nil {
		return nil, err
	}

	return toIngestToken(mutation.Result.IngestToken), err
}

func (i *IngestTokens) Update(repositoryName string, tokenName string, parserName string) (*IngestToken, error) {
	var mutation struct {
		Result struct {
			Repository Repository
		} `graphql:"assignIngestToken(repositoryName: $repositoryName, tokenName: $tokenName, parserName: $parserName)"`
	}

	variables := map[string]interface{}{
		"tokenName":      graphql.String(tokenName),
		"repositoryName": graphql.String(repositoryName),
		"parserName":     graphql.String(parserName),
	}

	err := i.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	return i.Get(repositoryName, tokenName)
}

func (i *IngestTokens) Remove(repositoryName string, tokenName string) error {
	var mutation struct {
		Result struct {
			Type string `graphql:"__typename"`
		} `graphql:"removeIngestToken(repositoryName: $repositoryName, name: $tokenName)"`
	}

	variables := map[string]interface{}{
		"tokenName":      graphql.String(tokenName),
		"repositoryName": graphql.String(repositoryName),
	}

	return i.client.Mutate(&mutation, variables)
}
