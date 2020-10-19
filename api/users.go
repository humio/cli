package api

import (
	"errors"
	"github.com/shurcooL/graphql"
)

type Users struct {
	client *Client
}

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

type UserChangeSet struct {
	IsRoot      *bool
	FullName    *string
	Company     *string
	CountryCode *string
	Picture     *string
	Email       *string
}

func (c *Client) Users() *Users { return &Users{client: c} }

func (u *Users) List() ([]User, error) {
	var q struct {
		Users []User `graphql:"users"`
	}

	graphqlErr := u.client.Query(&q, nil)

	return q.Users, graphqlErr
}

func (u *Users) Get(username string) (User, error) {
	var q struct {
		Users []User `graphql:"users(search: $username)"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
	}

	graphqlErr := u.client.Query(&q, variables)

	if graphqlErr != nil {
		return User{}, graphqlErr
	}

	for _, user := range q.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (u *Users) Update(username string, changeset UserChangeSet) (User, error) {
	var mutation struct {
		Result struct{ User User } `graphql:"updateUser(input: {username: $username, isRoot: $isRoot, fullName: $fullName, company: $company, countryCode: $countryCode, email: $email, picture: $picture})"`
	}

	graphqlErr := u.client.Mutate(&mutation, userChangesetToVars(username, changeset))

	return mutation.Result.User, graphqlErr
}

func (u *Users) Add(username string, changeset UserChangeSet) (User, error) {
	var mutation struct {
		Result struct{ User User } `graphql:"addUser(input: {username: $username, isRoot: $isRoot, fullName: $fullName, company: $company, countryCode: $countryCode, email: $email, picture: $picture})"`
	}

	graphqlErr := u.client.Mutate(&mutation, userChangesetToVars(username, changeset))

	return mutation.Result.User, graphqlErr
}

func (u *Users) Remove(username string) (User, error) {
	var mutation struct {
		Result struct {
			// We have to make a selection, so just take __typename
			User User
		} `graphql:"removeUser(input: {username: $username})"`
	}

	variables := map[string]interface{}{
		"username": graphql.String(username),
	}

	graphqlErr := u.client.Mutate(&mutation, variables)

	return mutation.Result.User, graphqlErr
}

func (u *Users) RotateUserApiTokenAndGet(userID string) (string, error) {
	var mutation struct {
		RotateUserApiTokenMutation struct {
			RotateUserApiToken struct {
				Token string
			} `graphql:"rotateUserApiToken"`
		} `graphql:"rotateUserApiTokenAndGet(input:{id:$id})"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(userID),
	}

	err := u.client.Mutate(&mutation, variables)

	if err != nil {
		return "", err
	}

	return mutation.RotateUserApiTokenMutation.RotateUserApiToken.Token, nil
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
