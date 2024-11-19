package api

import (
	"errors"

	graphql "github.com/cli/shurcooL-graphql"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Users struct {
	client *Client
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type User struct {
	ID          string
	Username    string
	FullName    string
	Email       string
	Company     string
	CountryCode string
	Picture     string
	IsRoot      bool
	CreatedAt   string
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type UserChangeSet struct {
	IsRoot      *bool
	FullName    *string
	Company     *string
	CountryCode *string
	Picture     *string
	Email       *string
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Users() *Users { return &Users{client: c} }

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (u *Users) List() ([]User, error) {
	var query struct {
		Users []User `graphql:"users"`
	}

	err := u.client.Query(&query, nil)
	return query.Users, err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (u *Users) Get(username string) (User, error) {
	var query struct {
		Users []User `graphql:"users(search: $username)"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
	}

	err := u.client.Query(&query, variables)
	if err != nil {
		return User{}, err
	}

	for _, user := range query.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (u *Users) Update(username string, changeset UserChangeSet) (User, error) {
	var mutation struct {
		Result struct{ User User } `graphql:"updateUser(input: {username: $username, company: $company, isRoot: $isRoot, fullName: $fullName, picture: $picture, email: $email, countryCode: $countryCode})"`
	}

	err := u.client.Mutate(&mutation, userChangesetToVars(username, changeset))
	return mutation.Result.User, err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (u *Users) Add(username string, changeset UserChangeSet) (User, error) {
	var mutation struct {
		Result struct {
			// We have to make a selection, so just take __typename
			Typename graphql.String `graphql:"__typename"`
		} `graphql:"addUserV2(input: {username: $username, company: $company, isRoot: $isRoot, fullName: $fullName, picture: $picture, email: $email, countryCode: $countryCode})"`
	}

	err := u.client.Mutate(&mutation, userChangesetToVars(username, changeset))
	if err != nil {
		return User{}, err
	}

	return u.Get(username)
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (u *Users) Remove(username string) (User, error) {
	var mutation struct {
		Result struct {
			User User
		} `graphql:"removeUser(input: {username: $username})"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
	}

	err := u.client.Mutate(&mutation, variables)
	return mutation.Result.User, err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (u *Users) RotateToken(userID string) (string, error) {
	var mutation struct {
		Token string `graphql:"rotateToken(input:{id:$id})"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(userID),
	}

	err := u.client.Mutate(&mutation, variables)
	if err != nil {
		return "", err
	}

	return mutation.Token, nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
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
