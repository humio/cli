package command

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/shurcooL/graphql"
	cli "gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
)

func ParserPush(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureURL(config)

	var content []byte
	var readErr error

	repoName := c.Args().Get(0)

	if c.IsSet("file") {
		filePath := c.String("file")
		content, readErr = getParserFromFile(filePath)
	} else if c.IsSet("url") {
		url := c.String("url")
		content, readErr = getUrlParser(url)
	} else {
		parserName := c.Args().Get(1)
		content, readErr = getGithubParser(parserName)
	}

	check(readErr)

	t := parserConfig{}

	err := yaml.Unmarshal(content, &t)
	check(err)

	client := newGraphQLClient(config)

	var m struct {
		CreateParser struct {
			Type string `graphql:"__typename"`
		} `graphql:"createParser(input: { name: $name, repositoryName: $repositoryName, testData: $testData, tagFields: $tagFields, sourceCode: $sourceCode, force: $force})"`
	}

	tagFields := make([]graphql.String, 0)

	var effectiveName string
	if c.IsSet("name") {
		effectiveName = c.String("name")
	} else {
		effectiveName = t.Name
	}

	log.Println("NAME:" + effectiveName)

	variables := map[string]interface{}{
		"name":           graphql.String(effectiveName),
		"sourceCode":     graphql.String(t.Script),
		"repositoryName": graphql.String(repoName),
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

func getParserFromFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func getGithubParser(parserName string) ([]byte, error) {
	url := "https://raw.githubusercontent.com/humio/community/master/parsers/" + parserName + ".yaml"
	return getUrlParser(url)
}

func getUrlParser(url string) ([]byte, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
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
