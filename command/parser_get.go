package command

import (
	"context"
	"io/ioutil"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
	cli "gopkg.in/urfave/cli.v2"
	yaml "gopkg.in/yaml.v2"
)

func ParserGet(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureRepo(config)
	ensureURL(config)

	parserName := c.Args().First()

	var outFilePath string
	if c.IsSet("out") {
		outFilePath = c.String("out")
	} else {
		outFilePath = parserName + ".yaml"
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)

	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient(config.URL+"graphql", httpClient)

	var q struct {
		Repository struct {
			Parser struct {
				Name       string
				SourceCode string
				TestData   []string
			} `graphql:"parser(name: $parserName)"`
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"parserName":     graphql.String(parserName),
		"repositoryName": graphql.String(config.Repo),
	}

	graphqlErr := client.Query(context.Background(), &q, variables)
	check(graphqlErr)

	yamlContent := parserConfig{
		Name:   q.Repository.Parser.Name,
		Tests:  Map(q.Repository.Parser.TestData, toTestCase),
		Script: q.Repository.Parser.SourceCode,
	}

	yamlData, yamlErr := yaml.Marshal(&yamlContent)
	check(yamlErr)
	ioutil.WriteFile(outFilePath, yamlData, 0644)

	return nil
}

func toTestCase(line string) testCase {
	return testCase{
		Input:  line,
		Output: map[string]string{},
	}
}
