package command

import (
	"context"
	"io/ioutil"
	"log"
	"strings"

	"github.com/shurcooL/graphql"
	cli "gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
)

func ParserAdd(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureRepo(config)
	ensureURL(config)

	filePath := c.Args().First()
	file, readErr := ioutil.ReadFile(filePath)
	check(readErr)

	t := parserConfig{}

	err := yaml.Unmarshal(file, &t)
	check(err)

	client := newGraphQLClient(config)

	var m struct {
		CreateParser struct {
			Type string `graphql:"__typename"`
		} `graphql:"createParser(input: { name: $name, repositoryName: $repositoryName, testData: $testData, tagFields: $tagFields, sourceCode: $sourceCode, force: $force})"`
	}

	tagFields := make([]graphql.String, 0)

	var effectiveName string
	if c.String("name") != "" {
		effectiveName = c.String("name")
	} else {
		effectiveName = t.Name
	}

	variables := map[string]interface{}{
		"name":           graphql.String(effectiveName),
		"sourceCode":     graphql.String(t.Script),
		"repositoryName": graphql.String(config.Repo),
		"testData":       testCasesToStrings(t),
		"tagFields":      tagFields,
		"force":          graphql.Boolean(true),
	}

	graphqlErr := client.Mutate(context.Background(), &m, variables)

	if graphqlErr != nil {
		log.Fatal(graphqlErr)
	}

	return nil
}

func testCasesToStrings(parser parserConfig) []graphql.String {

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
