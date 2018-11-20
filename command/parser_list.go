package command

import (
	"context"
	"fmt"
	"log"

	"github.com/ryanuber/columnize"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
	cli "gopkg.in/urfave/cli.v2"
)

func ParserList(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureRepo(config)
	ensureURL(config)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)

	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient(config.URL+"graphql", httpClient)

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

	graphqlErr := client.Query(context.Background(), &q, variables)

	if graphqlErr != nil {
		log.Fatal(graphqlErr)
	}

	var output []string
	output = append(output, "Name | Custom")
	for i := 0; i < len(q.Repository.Parsers); i++ {
		parser := q.Repository.Parsers[i]
		output = append(output, fmt.Sprintf("%v | %v", parser.Name, checkmark(!parser.IsBuiltIn)))
	}

	table := columnize.SimpleFormat(output)

	fmt.Println()
	fmt.Println(table)
	fmt.Println()

	return nil
}

func checkmark(v bool) string {
	if v {
		return "✓"
	}
	return ""
}
