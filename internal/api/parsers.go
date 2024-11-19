package api

import (
	"context"
	"fmt"

	"github.com/humio/cli/internal/api/humiographql"
)

const LogScaleVersionWithParserAPIv2 = "1.129.0"

type ParserTestEvent struct {
	RawString string `yaml:"rawString"`
}

type ParserTestCaseAssertions struct {
	OutputEventIndex int               `yaml:"outputEventIndex"`
	FieldsNotPresent []string          `yaml:"fieldsNotPresent"`
	FieldsHaveValues map[string]string `yaml:"fieldsHaveValues"`
}

type ParserTestCase struct {
	Event      ParserTestEvent            `yaml:"event"`
	Assertions []ParserTestCaseAssertions `yaml:"assertions"`
}

type Parser struct {
	ID                             string
	Name                           string
	Script                         string           `yaml:",flow"`
	TestCases                      []ParserTestCase `yaml:"testCases"`
	FieldsToTag                    []string         `yaml:"tagFields"`
	FieldsToBeRemovedBeforeParsing []string         `yaml:"fieldsToBeRemovedBeforeParsing"`
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
	resp, err := humiographql.ListParsers(context.Background(), p.client, repositoryName)
	if err != nil {
		return nil, err
	}

	respRepo := resp.GetRepository()
	respParsers := respRepo.GetParsers()
	parsers := make([]ParserListItem, len(respParsers))
	for idx, parser := range respParsers {
		parsers[idx] = ParserListItem{
			ID:        parser.GetId(),
			Name:      parser.GetName(),
			IsBuiltIn: parser.GetIsBuiltIn(),
		}
	}
	return parsers, nil
}

func (p *Parsers) Delete(repositoryName string, parserName string) error {
	status, getStatusErr := p.client.Status()
	if getStatusErr != nil {
		return getStatusErr
	}

	atLeast, versionParseErr := status.AtLeast(LogScaleVersionWithParserAPIv2)
	if versionParseErr != nil {
		return versionParseErr
	}
	parser, err := p.client.Parsers().Get(repositoryName, parserName)
	if err != nil {
		return err
	}
	if !atLeast {
		_, err = humiographql.LegacyDeleteParserByID(context.Background(), p.client, repositoryName, parser.ID)
		return err
	}

	_, err = humiographql.DeleteParserByID(context.Background(), p.client, repositoryName, parser.ID)
	return err
}

func (p *Parsers) Add(repositoryName string, newParser *Parser, allowOverwritingExistingParser bool) (*Parser, error) {
	if newParser == nil {
		return nil, fmt.Errorf("newFilterAlert must not be nil")
	}
	status, getStatusErr := p.client.Status()
	if getStatusErr != nil {
		return nil, getStatusErr
	}
	atLeast, versionParseErr := status.AtLeast(LogScaleVersionWithParserAPIv2)
	if versionParseErr != nil {
		return nil, versionParseErr
	}

	if !atLeast {
		testData := make([]string, len(newParser.TestCases))
		for i, testCase := range newParser.TestCases {
			testData[i] = testCase.Event.RawString
		}
		resp, err := humiographql.LegacyCreateParser(
			context.Background(),
			p.client,
			repositoryName,
			newParser.Name,
			testData,
			newParser.FieldsToTag,
			newParser.Script,
			allowOverwritingExistingParser,
		)
		if err != nil {
			return nil, err
		}

		respCreateParser := resp.GetCreateParser()
		respParser := respCreateParser.GetParser()
		respTestCases := respParser.GetTestCases()
		testCases := make([]ParserTestCase, len(respTestCases))
		for idx, testCase := range respTestCases {
			event := testCase.GetEvent()
			testCases[idx] = ParserTestCase{
				Event: ParserTestEvent{
					RawString: event.GetRawString(),
				},
				Assertions: nil,
			}
		}
		return &Parser{
			ID:                             respParser.GetId(),
			Name:                           respParser.GetName(),
			Script:                         respParser.GetScript(),
			TestCases:                      testCases,
			FieldsToTag:                    respParser.GetFieldsToTag(),
			FieldsToBeRemovedBeforeParsing: respParser.GetFieldsToBeRemovedBeforeParsing(),
		}, nil
	}

	testCasesInput := make([]humiographql.ParserTestCaseInput, len(newParser.TestCases))
	for j, pa := range newParser.TestCases {
		parserTestCaseAssertionsForOutputInput := make([]humiographql.ParserTestCaseAssertionsForOutputInput, len(pa.Assertions))
		for i := range pa.Assertions {
			fieldsHaveValuesInput := make([]humiographql.FieldHasValueInput, len(pa.Assertions[i].FieldsHaveValues))
			for field, value := range pa.Assertions[i].FieldsHaveValues {
				fieldsHaveValuesInput[i] = humiographql.FieldHasValueInput{
					FieldName:     field,
					ExpectedValue: value,
				}
			}
			parserTestCaseAssertionsForOutputInput[i] = humiographql.ParserTestCaseAssertionsForOutputInput{
				OutputEventIndex: pa.Assertions[i].OutputEventIndex,
				Assertions: humiographql.ParserTestCaseOutputAssertionsInput{
					FieldsNotPresent: pa.Assertions[i].FieldsNotPresent,
					FieldsHaveValues: fieldsHaveValuesInput,
				},
			}
		}
		testCasesInput[j] = humiographql.ParserTestCaseInput{
			Event:            humiographql.ParserTestEventInput{RawString: pa.Event.RawString},
			OutputAssertions: parserTestCaseAssertionsForOutputInput,
		}
	}
	resp, err := humiographql.CreateParser(
		context.Background(),
		p.client,
		repositoryName,
		newParser.Name,
		newParser.Script,
		testCasesInput,
		newParser.FieldsToTag,
		newParser.FieldsToBeRemovedBeforeParsing,
		allowOverwritingExistingParser,
	)
	if err != nil {
		return nil, err
	}

	respParser := resp.GetCreateParserV2()
	respTestCases := respParser.GetTestCases()
	testCases := make([]ParserTestCase, len(respTestCases))
	for idx, testCase := range respTestCases {
		testEvent := testCase.GetEvent()
		respOutputAssertions := testCase.GetOutputAssertions()
		assertions := make([]ParserTestCaseAssertions, len(respOutputAssertions))
		for j, assertion := range respOutputAssertions {
			respAssertion := assertion.GetAssertions()
			fieldHaveValue := respAssertion.GetFieldsHaveValues()
			fieldsHaveValues := map[string]string{}
			for _, field := range fieldHaveValue {
				fieldsHaveValues[field.GetFieldName()] = field.GetFieldName()
			}

			assertions[j] = ParserTestCaseAssertions{
				OutputEventIndex: assertion.GetOutputEventIndex(),
				FieldsNotPresent: respAssertion.FieldsNotPresent,
				FieldsHaveValues: nil,
			}
		}
		testCases[idx] = ParserTestCase{
			Event: ParserTestEvent{
				RawString: testEvent.GetRawString(),
			},
			Assertions: assertions,
		}
	}
	return &Parser{
		ID:                             respParser.GetId(),
		Name:                           respParser.GetName(),
		Script:                         respParser.GetScript(),
		TestCases:                      testCases,
		FieldsToTag:                    respParser.GetFieldsToTag(),
		FieldsToBeRemovedBeforeParsing: respParser.GetFieldsToBeRemovedBeforeParsing(),
	}, nil
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
		resp, err := humiographql.LegacyGetParser(
			context.Background(),
			p.client,
			repositoryName,
			parserName,
		)
		if err != nil {
			return nil, err
		}

		respRepository := resp.GetRepository()
		respParser := respRepository.GetParser()
		respTestCases := respParser.GetTestData()
		testCases := make([]ParserTestCase, len(respTestCases))
		for idx, testCase := range respTestCases {
			testCases[idx] = ParserTestCase{
				Event: ParserTestEvent{
					RawString: testCase,
				},
			}
		}
		return &Parser{
			ID:          respParser.GetId(),
			Name:        respParser.GetName(),
			Script:      respParser.GetSourceCode(),
			TestCases:   testCases,
			FieldsToTag: respParser.GetTagFields(),
		}, nil
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

	resp, err := humiographql.GetParserByID(
		context.Background(),
		p.client,
		repositoryName,
		parserID,
	)
	if err != nil {
		return nil, err
	}

	respRepository := resp.GetRepository()
	respParser := respRepository.GetParser()
	respTestCases := respParser.GetTestCases()
	testCases := make([]ParserTestCase, len(respTestCases))
	for idx, testCase := range respTestCases {
		testEvent := testCase.GetEvent()
		respOutputAssertions := testCase.GetOutputAssertions()
		assertions := make([]ParserTestCaseAssertions, len(respOutputAssertions))
		for j, assertion := range respOutputAssertions {
			respAssertion := assertion.GetAssertions()
			fieldHaveValue := respAssertion.GetFieldsHaveValues()
			fieldsHaveValues := map[string]string{}
			for _, field := range fieldHaveValue {
				fieldsHaveValues[field.GetFieldName()] = field.GetFieldName()
			}

			assertions[j] = ParserTestCaseAssertions{
				OutputEventIndex: assertion.GetOutputEventIndex(),
				FieldsNotPresent: respAssertion.FieldsNotPresent,
				FieldsHaveValues: nil,
			}
		}
		testCases[idx] = ParserTestCase{
			Event: ParserTestEvent{
				RawString: testEvent.GetRawString(),
			},
			Assertions: assertions,
		}
	}
	return &Parser{
		ID:                             respParser.GetId(),
		Name:                           respParser.GetName(),
		Script:                         respParser.GetScript(),
		TestCases:                      testCases,
		FieldsToTag:                    respParser.GetFieldsToTag(),
		FieldsToBeRemovedBeforeParsing: respParser.GetFieldsToBeRemovedBeforeParsing(),
	}, nil
}

func (p *Parsers) Export(repositoryName string, parserName string) (string, error) {
	resp, err := humiographql.GetParserYAMLByName(context.Background(), p.client, repositoryName, parserName)
	if err != nil {
		return "", err
	}
	respRepo := resp.GetRepository()
	return respRepo.GetParser().GetYamlTemplate(), err
}
