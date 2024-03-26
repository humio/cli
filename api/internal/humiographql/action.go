package humiographql

import graphql "github.com/cli/shurcooL-graphql"

type Action struct {
	Name graphql.String `graphql:"name"`
}
