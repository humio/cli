package humiographql

import (
	graphql "github.com/cli/shurcooL-graphql"
)

type UpdateParserScriptInput struct {
	Script graphql.String `json:"script"`
}

type ParserTestEventInput struct {
	RawString graphql.String `json:"rawString"`
}

type FieldHasValueInput struct {
	FieldName     graphql.String `json:"fieldName"`
	ExpectedValue graphql.String `json:"expectedValue"`
}

type ParserTestCaseOutputAssertionsInput struct {
	FieldsNotPresent []graphql.String     `json:"fieldsNotPresent"`
	FieldsHaveValues []FieldHasValueInput `json:"fieldsHaveValues"`
}

type ParserTestCaseAssertionsForOutputInput struct {
	OutputEventIndex graphql.Int                         `json:"outputEventIndex"`
	Assertions       ParserTestCaseOutputAssertionsInput `json:"assertions"`
}

type ParserTestCaseInput struct {
	Event            ParserTestEventInput                     `json:"event"`
	OutputAssertions []ParserTestCaseAssertionsForOutputInput `json:"outputAssertions"`
}

type ParserTestEvent struct {
	RawString graphql.String `graphql:"rawString"`
}

type FieldHasValue struct {
	FieldName     graphql.String `graphql:"fieldName"`
	ExpectedValue graphql.String `graphql:"expectedValue"`
}

type ParserTestCaseOutputAssertions struct {
	FieldsNotPresent []string        `graphql:"fieldsNotPresent"`
	FieldsHaveValues []FieldHasValue `graphql:"fieldsHaveValues"`
}

type ParserTestCaseAssertionsForOutput struct {
	OutputEventIndex graphql.Int                    `graphql:"outputEventIndex"`
	Assertions       ParserTestCaseOutputAssertions `graphql:"assertions"`
}

type ParserTestCase struct {
	Event            ParserTestEvent                     `graphql:"event"`
	OutputAssertions []ParserTestCaseAssertionsForOutput `graphql:"outputAssertions"`
}

type Parser struct {
	ID                             graphql.String   `graphql:"id"`
	Name                           graphql.String   `graphql:"name"`
	DisplayName                    graphql.String   `graphql:"displayName"`
	Description                    graphql.String   `graphql:"description"`
	IsBuiltIn                      graphql.Boolean  `graphql:"isBuiltIn"`
	Script                         graphql.String   `graphql:"script"`
	FieldsToTag                    []graphql.String `graphql:"fieldsToTag"`
	FieldsToBeRemovedBeforeParsing []graphql.String `graphql:"fieldsToBeRemovedBeforeParsing"`
	TestCases                      []ParserTestCase `graphql:"testCases"`
}

type CreateParserInputV2 struct {
	Name                           graphql.String        `json:"name"`
	Script                         graphql.String        `json:"script"`
	TestCases                      []ParserTestCaseInput `json:"testCases"`
	RepositoryName                 RepoOrViewName        `json:"repositoryName"`
	FieldsToTag                    []graphql.String      `json:"fieldsToTag"`
	FieldsToBeRemovedBeforeParsing []graphql.String      `json:"fieldsToBeRemovedBeforeParsing"`
	AllowOverwritingExistingParser graphql.Boolean       `json:"allowOverwritingExistingParser"`
}
