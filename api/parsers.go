package api

type ParserTestCase struct {
	Input  string
	Output map[string]string
}

type Parser struct {
	ID        string
	Name      string
	Tests     []string `yaml:",omitempty"`
	Example   string   `yaml:",omitempty"`
	Script    string   `yaml:",flow"`
	TagFields []string `yaml:",omitempty"`
}

type Parsers struct {
	client *Client
}

func (c *Client) Parsers() *Parsers { return &Parsers{client: c} }

type ParserListItem struct {
	ID        string
	Name      string
	IsBuiltIn bool
}

func (p *Parsers) List(repositoryName string) ([]ParserListItem, error) {
	var query struct {
		Repository struct {
			Parsers []ParserListItem
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": repositoryName,
	}

	var parsers []ParserListItem
	err := p.client.Query(&query, variables)
	if err == nil {
		parsers = query.Repository.Parsers
	}
	return parsers, err
}

func (p *Parsers) Remove(repositoryName string, parserName string) error {
	var mutation struct {
		RemoveParser struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"removeParser(input: { id: $id, repositoryName: $repositoryName })"`
	}

	parser, err := p.client.Parsers().Get(repositoryName, parserName)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"repositoryName": repositoryName,
		"id":             parser.ID,
	}

	return p.client.Mutate(&mutation, variables)
}

func (p *Parsers) Add(repositoryName string, parser *Parser, force bool) error {

	var mutation struct {
		CreateParser struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"createParser(input: { name: $name, repositoryName: $repositoryName, testData: $testData, tagFields: $tagFields, sourceCode: $sourceCode, force: $force})"`
	}

	tagFieldsGQL := make([]string, len(parser.TagFields))
	copy(tagFieldsGQL, parser.TagFields)

	testsGQL := make([]string, len(parser.Tests))
	copy(testsGQL, parser.Tests)

	variables := map[string]interface{}{
		"name":           parser.Name,
		"sourceCode":     parser.Script,
		"repositoryName": repositoryName,
		"testData":       testsGQL,
		"tagFields":      tagFieldsGQL,
		"force":          force,
	}

	return p.client.Mutate(&mutation, variables)
}

func (p *Parsers) Get(repositoryName string, parserName string) (*Parser, error) {
	var query struct {
		Repository struct {
			Parser *struct {
				ID         string
				Name       string
				SourceCode string
				TestData   []string
				TagFields  []string
			} `graphql:"parser(name: $parserName)"`
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"parserName":     parserName,
		"repositoryName": repositoryName,
	}

	err := p.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	if query.Repository.Parser == nil {
		return nil, ParserNotFound(parserName)
	}

	parser := Parser{
		ID:        query.Repository.Parser.ID,
		Name:      query.Repository.Parser.Name,
		Tests:     query.Repository.Parser.TestData,
		Script:    query.Repository.Parser.SourceCode,
		TagFields: query.Repository.Parser.TagFields,
	}

	return &parser, nil
}

func (p *Parsers) Export(repositoryName string, parserName string) (string, error) {

	var query struct {
		Repository struct {
			Parser *struct {
				Name         string
				YamlTemplate string
			} `graphql:"parser(name: $parserName)"`
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"parserName":     parserName,
		"repositoryName": repositoryName,
	}

	err := p.client.Query(&query, variables)

	if err != nil {
		return "", err
	}

	if query.Repository.Parser == nil {
		return "", ParserNotFound(parserName)
	}

	return query.Repository.Parser.YamlTemplate, nil
}
