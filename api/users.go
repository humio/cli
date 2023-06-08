package api

import (
	"errors"
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
	var query struct {
		Users []User `graphql:"users"`
	}

	err := u.client.Query(&query, nil)
	return query.Users, err
}

func (u *Users) Get(username string) (User, error) {
	var query struct {
		Users []User `graphql:"users(search: $username)"`
	}

	variables := map[string]interface{}{
		"username": username,
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

func (u *Users) Update(username string, changeset UserChangeSet) (User, error) {
	var mutation struct {
		Result struct{ User User } `graphql:"updateUser(input: {username: $username, company: $company, isRoot: $isRoot, fullName: $fullName, picture: $picture, email: $email, countryCode: $countryCode})"`
	}

	err := u.client.Mutate(&mutation, userChangesetToVars(username, changeset))
	return mutation.Result.User, err
}

func (u *Users) Add(username string, changeset UserChangeSet) (User, error) {
	var mutation struct {
		Result struct {
			// We have to make a selection, so just take __typename
			Typename string `graphql:"__typename"`
		} `graphql:"addUserV2(input: {username: $username, company: $company, isRoot: $isRoot, fullName: $fullName, picture: $picture, email: $email, countryCode: $countryCode})"`
	}

	err := u.client.Mutate(&mutation, userChangesetToVars(username, changeset))
	if err != nil {
		return User{}, err
	}

	return u.Get(username)
}

func (u *Users) Remove(username string) (User, error) {
	var mutation struct {
		Result struct {
			User User
		} `graphql:"removeUser(input: {username: $username})"`
	}

	variables := map[string]interface{}{
		"username": username,
	}

	err := u.client.Mutate(&mutation, variables)
	return mutation.Result.User, err
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
		"id": userID,
	}

	err := u.client.Mutate(&mutation, variables)
	if err != nil {
		return "", err
	}

	return mutation.RotateUserApiTokenMutation.RotateUserApiToken.Token, nil
}

func userChangesetToVars(username string, changeset UserChangeSet) map[string]interface{} {
	return map[string]interface{}{
		"username":    username,
		"isRoot":      changeset.IsRoot,
		"fullName":    changeset.FullName,
		"company":     changeset.Company,
		"countryCode": changeset.CountryCode,
		"email":       changeset.Email,
		"picture":     changeset.Picture,
	}
}
