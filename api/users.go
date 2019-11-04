package api

import "github.com/shurcooL/graphql"

type Users struct {
	client *Client
}

type Role struct {
	Name string
}

type User struct {
	Username    string
	FullName    string
	Email       string
	Company     string
	CountryCode string
	Picture     string
	IsRoot      bool
	CreatedAt   string
	Roles       []Role
}

type UserChangeSet struct {
	IsRoot      *bool
	FullName    *string
	Company     *string
	CountryCode *string
	Picture     *string
	Email       *string
}

func (c *Client) Users() *Users { return &Users{client: c} }

func (c *Users) List() ([]User, error) {
	var q struct {
		Users []User `graphql:"accounts"`
	}

	graphqlErr := c.client.Query(&q, nil)

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
		Result struct{ User User } `graphql:"updateUser(input: {username: $username, isRoot: $isRoot, fullName: $fullName, company: $company, countryCode: $countryCode, email: $email, picture: $picture})"`
	}

	graphqlErr := c.client.Mutate(&mutation, userChangesetToVars(username, changeset))

	return mutation.Result.User, graphqlErr
}

func (c *Users) Add(username string, changeset UserChangeSet) (User, error) {
	var mutation struct {
		Result struct{ User User } `graphql:"addUser(input: {username: $username, isRoot: $isRoot, fullName: $fullName, company: $company, countryCode: $countryCode, email: $email, picture: $picture})"`
	}

	graphqlErr := c.client.Mutate(&mutation, userChangesetToVars(username, changeset))

	return mutation.Result.User, graphqlErr
}

func (c *Users) Remove(username string) (User, error) {
	var mutation struct {
		Result struct {
			// We have to make a selection, so just take __typename
			User User
		} `graphql:"removeUser(input: {username: $username})"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
	}

	graphqlErr := c.client.Mutate(&mutation, variables)

	return mutation.Result.User, graphqlErr
}

func userChangesetToVars(username string, changeset UserChangeSet) map[string]interface{} {
	return map[string]interface{}{
		"username":    graphql.String(username),
		"isRoot":      optBoolArg(changeset.IsRoot),
		"fullName":    optStringArg(changeset.FullName),
		"company":     optStringArg(changeset.Company),
		"countryCode": optStringArg(changeset.CountryCode),
		"email":       optStringArg(changeset.Email),
		"picture":     optStringArg(changeset.Picture),
	}
}
