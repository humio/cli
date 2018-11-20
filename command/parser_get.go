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

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)

	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient(config.URL+"graphql", httpClient)

	var q struct {
		Repository struct {
			Parser struct {
				SourceCode string
			} `graphql:"parser(name: $parserName)"`
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"parserName":     graphql.String(parserName),
		"repositoryName": graphql.String(config.Repo),
	}

	graphqlErr := client.Query(context.Background(), &q, variables)
	check(graphqlErr)

	d, yamlErr := yaml.Marshal(&q)
	check(yamlErr)

	ioutil.WriteFile(parserName+".yaml", d, 0644)

	return nil
}
