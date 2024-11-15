package api

import (
	"context"

	"github.com/humio/cli/internal/api/humiographql"
)

type IngestTokens struct {
	client *Client
}

type IngestToken struct {
	Name           string
	Token          string
	AssignedParser string
}

func (c *Client) IngestTokens() *IngestTokens { return &IngestTokens{client: c} }

func (i *IngestTokens) List(repositoryName string) ([]IngestToken, error) {
	resp, err := humiographql.ListIngestTokens(context.Background(), i.client, repositoryName)
	if err != nil {
		return nil, err
	}
	respRepo := resp.GetRepository()
	respRepoTokens := respRepo.GetIngestTokens()
	tokens := make([]IngestToken, len(respRepoTokens))
	for idx, token := range respRepoTokens {
		respParser := token.GetParser()
		assignedParser := ""
		if respParser != nil {
			assignedParser = respParser.GetName()
		}
		tokens[idx] = IngestToken{
			Name:           token.GetName(),
			Token:          token.GetToken(),
			AssignedParser: assignedParser,
		}
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

	return nil, IngestTokenNotFound(tokenName)
}

func (i *IngestTokens) Add(repositoryName string, tokenName string, parserName string) (*IngestToken, error) {
	var parserNamePtr *string
	if parserName != "" {
		parserNamePtr = &parserName
	}
	resp, err := humiographql.AddIngestToken(context.Background(), i.client, repositoryName, tokenName, parserNamePtr)
	if err != nil {
		return nil, err
	}

	respToken := resp.GetAddIngestTokenV3()
	assignedParser := ""
	respParser := respToken.GetParser()
	if respParser != nil {
		assignedParser = respParser.GetName()
	}
	return &IngestToken{
		Name:           respToken.GetName(),
		Token:          respToken.GetToken(),
		AssignedParser: assignedParser,
	}, nil
}

func (i *IngestTokens) Update(repositoryName string, tokenName string, parserName string) (*IngestToken, error) {
	if parserName == "" {
		_, err := humiographql.UnassignParserToIngestToken(context.Background(), i.client, repositoryName, tokenName)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := humiographql.AssignParserToIngestToken(context.Background(), i.client, repositoryName, tokenName, parserName)
		if err != nil {
			return nil, err
		}
	}

	return i.Get(repositoryName, tokenName)
}

func (i *IngestTokens) Remove(repositoryName string, tokenName string) error {
	_, err := humiographql.RemoveIngestToken(context.Background(), i.client, repositoryName, tokenName)
	return err
}
