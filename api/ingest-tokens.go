package api

import (
	"encoding/json"
	"io/ioutil"
	"log"

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

func (p *IngestTokens) List(reposistoryName string) ([]IngestToken, error) {
	url := p.client.Address() + "api/v1/repositories/" + reposistoryName + "/ingesttokens"

	resp, clientErr := p.client.httpGET(url)
	defer resp.Body.Close()

	if clientErr != nil {
		log.Fatal(clientErr)
	}

	if resp.StatusCode >= 400 {
		log.Fatal(resp)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	tokens := make([]IngestToken, 0)
	jsonErr := json.Unmarshal(body, &tokens)
	return tokens, jsonErr
}

func (p *IngestTokens) Add(repo string, name string, parserName *string) (*IngestToken, error) {
	var mutation struct {
		Result struct {
			IngestToken struct {
				Name   string
				Token  string
				Parser *struct {
					Name string
				}
			}
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

	var result IngestToken
	if err == nil {
		var parser string
		if mutation.Result.IngestToken.Parser != nil {
			parser = mutation.Result.IngestToken.Parser.Name
		}

		result = IngestToken{
			Name:           mutation.Result.IngestToken.Name,
			Token:          mutation.Result.IngestToken.Token,
			AssignedParser: parser,
		}
	}

	return &result, err
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
