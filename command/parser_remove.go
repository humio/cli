package command

import (
	"context"
	"log"

	"github.com/shurcooL/graphql"
	cli "gopkg.in/urfave/cli.v2"
)

func ParserRemove(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureRepo(config)
	ensureURL(config)

	parserName := c.Args().First()

	client := newGraphQLClient(config)

	var m struct {
		CreateParser struct {
			Type string `graphql:"__typename"`
		} `graphql:"removeParser(input: { name: $name, repositoryName: $repositoryName })"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(config.Repo),
		"name":           graphql.String(parserName),
	}

	graphqlErr := client.Mutate(context.Background(), &m, variables)

	if graphqlErr != nil {
		log.Fatal(graphqlErr)
	}

	return nil
}
