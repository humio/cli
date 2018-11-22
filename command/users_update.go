package command

import (
	"context"

	"github.com/shurcooL/graphql"
	cli "gopkg.in/urfave/cli.v2"
)

func UpdateUser(c *cli.Context) error {
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
		"isRoot":   optBoolFlag(c, "root"),
	}

	graphqlErr := client.Mutate(context.Background(), &m, variables)
	check(graphqlErr)

	UsersShow(c)

	return nil
}

func optBoolFlag(c *cli.Context, flag string) *graphql.Boolean {
	var isRootOpt *graphql.Boolean
	if c.IsSet(flag) {
		isRootOpt = graphql.NewBoolean(graphql.Boolean(c.Bool(flag)))
	}
	return isRootOpt
}
