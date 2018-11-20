package command

import (
	"context"
	"strings"

	"github.com/shurcooL/graphql"
	cli "gopkg.in/urfave/cli.v2"
)

func UsersShow(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureURL(config)

	username := c.Args().First()

	var q struct {
		Account struct {
			Username  string
			FullName  string
			IsRoot    bool
			CreatedAt string
		} `graphql:"account(username: $username)"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
	}

	graphqlErr := newGraphQLClient(config).Query(context.Background(), &q, variables)
	check(graphqlErr)

	userData := []string{q.Account.Username, q.Account.FullName, q.Account.CreatedAt, yesNo(q.Account.IsRoot)}

	printTable([]string{
		"Username | Name | Created At | Is Root",
		strings.Join(userData, "|"),
	})

	return nil
}
