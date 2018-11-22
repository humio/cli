package api

import (
	"strings"

	"github.com/shurcooL/graphql"
)

type ParserTestCase struct {
	Input  string
	Output map[string]string
}

type Parser struct {
	Name        string
	Description string           `yaml:",omitempty"`
	Tests       []ParserTestCase `yaml:",omitempty"`
	Example     string           `yaml:",omitempty"`
	Script      string           `yaml:",flow"`
	TagFields   []string         `yaml:",omitempty"`
}

type Parsers struct {
	client *Client
}

func (c *Client) Parsers() *Parsers { return &Parsers{client: c} }

type ParserListItem struct {
	Name      string
	IsBuiltIn bool
}

func (p *Parsers) List(reposistoryName string) ([]ParserListItem, error) {
	var q struct {
		Repository struct {
			Parsers []ParserListItem
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(reposistoryName),
	}

	graphqlErr := p.client.Query(&q, variables)

	var parsers []ParserListItem
	if graphqlErr == nil {
		parsers = q.Repository.Parsers
	}

	return parsers, graphqlErr
}

func (p *Parsers) Remove(reposistoryName string, parserName string) error {
	var mutation struct {
		RemoveParser struct {
			Type string `graphql:"__typename"`
		} `graphql:"removeParser(input: { name: $name, repositoryName: $repositoryName })"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(reposistoryName),
		"name":           graphql.String(parserName),
	}

	return p.client.Mutate(&mutation, variables)
}

func (p *Parsers) Add(reposistoryName string, parser *Parser, force bool) error {

	var mutation struct {
		CreateParser struct {
			Type string `graphql:"__typename"`
		} `graphql:"createParser(input: { name: $name, repositoryName: $repositoryName, testData: $testData, tagFields: $tagFields, sourceCode: $sourceCode, force: $force})"`
	}

	tagFieldsGQL := make([]graphql.String, len(parser.TagFields))

	for i, field := range parser.TagFields {
		tagFieldsGQL[i] = graphql.String(field)
	}

	variables := map[string]interface{}{
		"name":           graphql.String(parser.Name),
		"sourceCode":     graphql.String(parser.Script),
		"repositoryName": graphql.String(reposistoryName),
		"testData":       testCasesToStrings(parser),
		"tagFields":      tagFieldsGQL,
		"force":          graphql.Boolean(force),
	}

	return p.client.Mutate(&mutation, variables)
}

func testCasesToStrings(parser *Parser) []graphql.String {

	lines := strings.Split(parser.Example, "\n")

	result := make([]graphql.String, 0)
	for _, item := range parser.Tests {
		result = append(result, graphql.String(item.Input))
	}

	for i, item := range lines {
		if i != len(lines)-1 {
			result = append(result, graphql.String(item))
		}
	}

	return result
}
