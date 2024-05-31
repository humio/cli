package api

import (
	"fmt"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/humio/cli/api/internal/humiographql"
)

const LogScaleVersionWithParserAPIv2 = "1.129.0"

type ParserTestEvent struct {
	RawString string `json:"rawString" yaml:"rawString"`
}

type ParserTestCaseAssertions struct {
	OutputEventIndex int               `json:"outputEventIndex" yaml:"outputEventIndex"`
	FieldsNotPresent []string          `json:"fieldsNotPresent" yaml:"fieldsNotPresent"`
	FieldsHaveValues map[string]string `json:"fieldsHaveValues" yaml:"fieldsHaveValues"`
}

type ParserTestCase struct {
	Event      ParserTestEvent            `json:"event"      yaml:"event"`
	Assertions []ParserTestCaseAssertions `json:"assertions" yaml:"assertions"`
}

type Parser struct {
	ID                             string
	Name                           string
	Script                         string           `json:"script"                                   yaml:",flow"`
	TestCases                      []ParserTestCase `json:"testCases"                                yaml:"testCases"`
	FieldsToTag                    []string         `json:"tagFields"                                yaml:"tagFields"`
	FieldsToBeRemovedBeforeParsing []string         `json:"fieldsToBeRemovedBeforeParsing,omitempty" yaml:"fieldsToBeRemovedBeforeParsing"`
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
		"repositoryName": graphql.String(repositoryName),
	}

	var parsers []ParserListItem
	err := p.client.Query(&query, variables)
	if err == nil {
		parsers = query.Repository.Parsers
	}
	return parsers, err
}

func (p *Parsers) Delete(repositoryName string, parserName string) error {
	status, err := p.client.Status()
	if err != nil {
		return err
	}

	atLeast, err := status.AtLeast(LogScaleVersionWithParserAPIv2)
	if err != nil {
		return err
	}
	if !atLeast {
		var mutation struct {
			RemoveParser struct {
				// We have to make a selection, so just take __typename
				Typename graphql.String `graphql:"__typename"`
			} `graphql:"removeParser(input: { id: $id, repositoryName: $repositoryName })"`
		}

		parser, err := p.client.Parsers().Get(repositoryName, parserName)
		if err != nil {
			return err
		}

		variables := map[string]interface{}{
			"repositoryName": graphql.String(repositoryName),
			"id":             graphql.String(parser.ID),
		}

		return p.client.Mutate(&mutation, variables)
	}

	var mutation struct {
		DeleteParser struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"deleteParser(input: { id: $id, repositoryName: $repositoryName })"`
	}

	parser, err := p.client.Parsers().Get(repositoryName, parserName)
	if err != nil {
		return err
	}

	variables := map[string]interface{}{
		"repositoryName": humiographql.RepoOrViewName(repositoryName),
		"id":             graphql.String(parser.ID),
	}

	return p.client.Mutate(&mutation, variables)
}

func (p *Parsers) Add(repositoryName string, newParser *Parser, allowOverwritingExistingParser bool) (*Parser, error) {
	if newParser == nil {
		return nil, fmt.Errorf("newFilterAlert must not be nil")
	}

	status, err := p.client.Status()
	if err != nil {
		return nil, err
	}

	atLeast, err := status.AtLeast(LogScaleVersionWithParserAPIv2)
	if err != nil {
		return nil, err
	}
	if !atLeast {
		var mutation struct {
			CreateParser struct {
				// We have to make a selection, so just take __typename
				Parser humiographql.Parser `graphql:"parser"`
			} `graphql:"createParser(input: { name: $name, repositoryName: $repositoryName, testData: $testData, tagFields: $tagFields, sourceCode: $sourceCode, force: $force})"`
		}

		testDataGQL := make([]graphql.String, len(newParser.TestCases))
		for i := range newParser.TestCases {
			testDataGQL[i] = graphql.String(newParser.TestCases[i].Event.RawString)
		}
		tagFieldsGQL := make([]graphql.String, len(newParser.FieldsToTag))
		for i := range newParser.FieldsToTag {
			tagFieldsGQL[i] = graphql.String(newParser.FieldsToTag[i])
		}

		variables := map[string]interface{}{
			"name":           graphql.String(newParser.Name),
			"sourceCode":     graphql.String(newParser.Script),
			"repositoryName": graphql.String(repositoryName),
			"testData":       testDataGQL,
			"tagFields":      tagFieldsGQL,
			"force":          graphql.Boolean(allowOverwritingExistingParser),
		}

		err = p.client.Mutate(&mutation, variables)
		if err != nil {
			return nil, err
		}

		parser := mapHumioGraphqlParserToParser(mutation.CreateParser.Parser)

		return &parser, nil
	}

	var mutation struct {
		humiographql.Parser `graphql:"createParserV2(input: $input)"`
	}

	fieldsToTagGQL := make([]graphql.String, len(newParser.FieldsToTag))
	for i, field := range newParser.FieldsToTag {
		fieldsToTagGQL[i] = graphql.String(field)
	}
	fieldsToBeRemovedBeforeParsingGQL := make([]graphql.String, len(newParser.FieldsToBeRemovedBeforeParsing))
	for i, field := range newParser.FieldsToBeRemovedBeforeParsing {
		fieldsToBeRemovedBeforeParsingGQL[i] = graphql.String(field)
	}
	testCasesGQL := make([]humiographql.ParserTestCaseInput, len(newParser.TestCases))
	for i := range newParser.TestCases {
		testCasesGQL[i] = mapParserTestCaseToInput(newParser.TestCases[i])
	}

	createParser := humiographql.CreateParserInputV2{
		Name:                           graphql.String(newParser.Name),
		Script:                         graphql.String(newParser.Script),
		TestCases:                      testCasesGQL,
		RepositoryName:                 humiographql.RepoOrViewName(repositoryName),
		FieldsToTag:                    fieldsToTagGQL,
		FieldsToBeRemovedBeforeParsing: fieldsToBeRemovedBeforeParsingGQL,
		AllowOverwritingExistingParser: graphql.Boolean(allowOverwritingExistingParser),
	}

	variables := map[string]interface{}{
		"input": createParser,
	}

	err = p.client.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	parser := mapHumioGraphqlParserToParser(mutation.Parser)

	return &parser, nil
}

func mapParserTestCaseToInput(p ParserTestCase) humiographql.ParserTestCaseInput {
	parserTestCaseAssertionsForOutputInput := make([]humiographql.ParserTestCaseAssertionsForOutputInput, len(p.Assertions))
	for i := range p.Assertions {
		fieldsNotPresent := make([]graphql.String, len(p.Assertions[i].FieldsNotPresent))
		for i := range p.Assertions[i].FieldsNotPresent {
			fieldsNotPresent[i] = graphql.String(p.Assertions[i].FieldsNotPresent[i])
		}
		fieldsHaveValuesInput := make([]humiographql.FieldHasValueInput, len(p.Assertions[i].FieldsHaveValues))
		for field, value := range p.Assertions[i].FieldsHaveValues {
			fieldsHaveValuesInput[i] = humiographql.FieldHasValueInput{
				FieldName:     graphql.String(field),
				ExpectedValue: graphql.String(value),
			}
		}
		parserTestCaseAssertionsForOutputInput[i] = humiographql.ParserTestCaseAssertionsForOutputInput{
			OutputEventIndex: graphql.Int(p.Assertions[i].OutputEventIndex),
			Assertions: humiographql.ParserTestCaseOutputAssertionsInput{
				FieldsNotPresent: fieldsNotPresent,
				FieldsHaveValues: fieldsHaveValuesInput,
			},
		}
	}
	return humiographql.ParserTestCaseInput{
		Event:            humiographql.ParserTestEventInput{RawString: graphql.String(p.Event.RawString)},
		OutputAssertions: parserTestCaseAssertionsForOutputInput,
	}
}

func (p *Parsers) Get(repositoryName string, parserName string) (*Parser, error) {
	status, err := p.client.Status()
	if err != nil {
		return nil, err
	}

	atLeast, err := status.AtLeast(LogScaleVersionWithParserAPIv2)
	if err != nil {
		return nil, err
	}
	if !atLeast {
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
			"parserName":     graphql.String(parserName),
			"repositoryName": graphql.String(repositoryName),
		}

		err := p.client.Query(&query, variables)
		if err != nil {
			return nil, err
		}

		if query.Repository.Parser == nil {
			return nil, ParserNotFound(parserName)
		}

		parser := Parser{
			ID:          query.Repository.Parser.ID,
			Name:        query.Repository.Parser.Name,
			Script:      query.Repository.Parser.SourceCode,
			FieldsToTag: query.Repository.Parser.TagFields,
		}
		parser.TestCases = make([]ParserTestCase, len(query.Repository.Parser.TestData))
		for i := range query.Repository.Parser.TestData {
			parser.TestCases[i] = ParserTestCase{
				Event: ParserTestEvent{RawString: query.Repository.Parser.TestData[i]},
			}
		}

		return &parser, nil
	}

	parserList, err := p.List(repositoryName)
	if err != nil {
		return nil, err
	}
	parserID := ""
	for i := range parserList {
		if parserList[i].Name == parserName {
			parserID = parserList[i].ID
			break
		}
	}
	if parserID == "" {
		return nil, ParserNotFound(parserName)
	}

	var query struct {
		Repository struct {
			Parser *humiographql.Parser `graphql:"parser(id: $parserID)"`
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"parserID":       graphql.String(parserID),
		"repositoryName": graphql.String(repositoryName),
	}

	err = p.client.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	parser := mapHumioGraphqlParserToParser(*query.Repository.Parser)

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
		"parserName":     graphql.String(parserName),
		"repositoryName": graphql.String(repositoryName),
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

func mapHumioGraphqlParserToParser(input humiographql.Parser) Parser {
	var fieldsToTag = make([]string, len(input.FieldsToTag))
	for i := range input.FieldsToTag {
		fieldsToTag[i] = string(input.FieldsToTag[i])
	}

	var fieldsToBeRemovedBeforeParsing = make([]string, len(input.FieldsToBeRemovedBeforeParsing))
	for i := range input.FieldsToBeRemovedBeforeParsing {
		fieldsToBeRemovedBeforeParsing[i] = string(input.FieldsToBeRemovedBeforeParsing[i])
	}

	var testCases = make([]ParserTestCase, len(input.TestCases))
	for i := range input.TestCases {
		var assertions = make([]ParserTestCaseAssertions, len(input.TestCases[i].OutputAssertions))
		for j := range input.TestCases[i].OutputAssertions {
			var fieldsHaveValues = make(map[string]string, len(input.TestCases[i].OutputAssertions[j].Assertions.FieldsHaveValues))
			for k := range input.TestCases[i].OutputAssertions[j].Assertions.FieldsHaveValues {
				fieldName := string(input.TestCases[i].OutputAssertions[j].Assertions.FieldsHaveValues[k].FieldName)
				expectedValue := string(input.TestCases[i].OutputAssertions[j].Assertions.FieldsHaveValues[k].ExpectedValue)
				fieldsHaveValues[fieldName] = expectedValue
			}

			assertions[j] = ParserTestCaseAssertions{
				OutputEventIndex: int(input.TestCases[i].OutputAssertions[j].OutputEventIndex),
				FieldsNotPresent: input.TestCases[i].OutputAssertions[j].Assertions.FieldsNotPresent,
				FieldsHaveValues: fieldsHaveValues,
			}
		}

		testCases[i] = ParserTestCase{
			Event:      ParserTestEvent{RawString: string(input.TestCases[i].Event.RawString)},
			Assertions: assertions,
		}
	}

	return Parser{
		ID:                             string(input.ID),
		Name:                           string(input.Name),
		Script:                         string(input.Script),
		TestCases:                      testCases,
		FieldsToTag:                    fieldsToTag,
		FieldsToBeRemovedBeforeParsing: fieldsToBeRemovedBeforeParsing,
	}
}
