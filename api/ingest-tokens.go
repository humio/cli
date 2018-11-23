package api

import (
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

func (p *IngestTokens) List(repo string) ([]IngestToken, error) {
	var query struct {
		Result struct {
			IngestTokens []ingestTokenData
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(repo),
	}

	err := p.client.Query(&query, variables)

	if err != nil {
		return nil, err
	}

	tokens := make([]IngestToken, len(query.Result.IngestTokens))
	for i, tokenData := range query.Result.IngestTokens {
		tokens[i] = *toIngestToken(tokenData)
	}

	return tokens, nil
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

func (p *IngestTokens) Add(repo string, name string, parserName *string) (*IngestToken, error) {
	var mutation struct {
		Result struct {
			IngestToken ingestTokenData
		} `graphql:"addIngestToken(repositoryName: $repositoryName, name: $name, parser: $parser)"`
	}

	var parserNameVar graphql.String
	if parserName != nil {
		parserNameVar = graphql.String(*parserName)
	}
	variables := map[string]interface{}{
		"name":           graphql.String(name),
		"repositoryName": graphql.String(repo),
		"parser":         parserNameVar,
	}

	err := p.client.Mutate(&mutation, variables)

	if err != nil {
		return nil, err
	}

	return toIngestToken(mutation.Result.IngestToken), err
}

func (p *IngestTokens) Remove(repo string, tokenName string) error {
	var mutation struct {
		Result struct {
			Type string `graphql:"__typename"`
		} `graphql:"removeIngestToken(repositoryName: $repositoryName, name: $name)"`
	}

	variables := map[string]interface{}{
		"name":           graphql.String(tokenName),
		"repositoryName": graphql.String(repo),
	}

	return p.client.Mutate(&mutation, variables)
}
