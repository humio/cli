package command

import (
	"context"

	cli "gopkg.in/urfave/cli.v2"
)

func UsersList(c *cli.Context) error {
	config, _ := getServerConfig(c)

	ensureToken(config)
	ensureURL(config)

	var q struct {
		Accounts []simpleAccount `graphql:"accounts"`
	}

	variables := map[string]interface{}{}

	graphqlErr := newGraphQLClient(config).Query(context.Background(), &q, variables)
	check(graphqlErr)

	rows := make([]string, len(q.Accounts))
	for i, account := range q.Accounts {
		rows[i] = formatSimpleAccount(account)
	}

	printTable(append([]string{
		"Username | Name | Root | Created"},
		rows...,
	))

	return nil
}
