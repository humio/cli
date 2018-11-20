package command

import (
	"context"
	"fmt"

	"github.com/shurcooL/graphql"
	cli "gopkg.in/urfave/cli.v2"
)

func ParserList(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureRepo(config)
	ensureURL(config)

	var q struct {
		Repository struct {
			Parsers []struct {
				Name      string
				IsBuiltIn bool
			}
		} `graphql:"repository(name: $repositoryName)"`
	}

	variables := map[string]interface{}{
		"repositoryName": graphql.String(config.Repo),
	}

	client := newGraphQLClient(config)
	graphqlErr := client.Query(context.Background(), &q, variables)
	check(graphqlErr)

	var output []string
	output = append(output, "Name | Custom")
	for i := 0; i < len(q.Repository.Parsers); i++ {
		parser := q.Repository.Parsers[i]
		output = append(output, fmt.Sprintf("%v | %v", parser.Name, checkmark(!parser.IsBuiltIn)))
	}

	printTable(output)

	return nil
}

func checkmark(v bool) string {
	if v {
		return "✓"
	}
	return ""
}
