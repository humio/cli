package api

import "github.com/shurcooL/graphql"

type Users struct {
	client *Client
}

type User struct {
	Username  string
	FullName  string
	IsRoot    bool
	CreatedAt string
}

type UserChangeSet struct {
	IsRoot *bool
}

func (c *Client) Users() *Users { return &Users{client: c} }

func (c *Users) List() ([]User, error) {
	var q struct {
		Users []User `graphql:"accounts"`
	}

	variables := map[string]interface{}{}

	graphqlErr := c.client.Query(&q, variables)

	return q.Users, graphqlErr
}

func (c *Users) Get(username string) (User, error) {
	var q struct {
		User User `graphql:"account(username: $username)"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
	}

	graphqlErr := c.client.Query(&q, variables)

	return q.User, graphqlErr
}

func (c *Users) Update(username string, changeset UserChangeSet) (User, error) {
	var mutation struct {
		Result struct{ User User } `graphql:"updateUser(input: {username: $username, isRoot: $isRoot})"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
		"isRoot":   optBoolArg(changeset.IsRoot),
	}

	graphqlErr := c.client.Mutate(&mutation, variables)

	return mutation.Result.User, graphqlErr
}

func (c *Users) Remove(username string, changeset UserChangeSet) error {
	var mutation struct {
		Result struct {
			// We have to make a selection, so just take __typename
			Type string `graphql:"__typename"`
		} `graphql:"removeUser(username: $username)"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
	}

	return c.client.Mutate(&mutation, variables)
}
