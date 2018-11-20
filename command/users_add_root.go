package command

import (
	"context"
	"log"

	"github.com/shurcooL/graphql"
	cli "gopkg.in/urfave/cli.v2"
)

func UsersAddRoot(c *cli.Context) error {
	return updateUser(c, true)
}

func UsersRemoveRoot(c *cli.Context) error {
	return updateUser(c, false)
}

func updateUser(c *cli.Context, isRoot bool) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureURL(config)

	username := c.Args().First()

	client := newGraphQLClient(config)

	var m struct {
		UpdateUser struct {
			Type string `graphql:"__typename"`
		} `graphql:"updateUser(input: { username: $username, isRoot: $isRoot })"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
		"isRoot":   graphql.Boolean(isRoot),
	}

	graphqlErr := client.Mutate(context.Background(), &m, variables)

	if graphqlErr != nil {
		log.Fatal(graphqlErr)
	} else if isRoot {
		log.Println(username + " now has root access to " + config.URL)
	} else {
		log.Println(username + " no longer has root access to " + config.URL)
	}

	return nil
}
