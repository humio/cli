package humiographql

import graphql "github.com/cli/shurcooL-graphql"

type User struct {
	ID graphql.String `graphql:"id"`
}
