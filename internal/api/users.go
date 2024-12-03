package api

import (
	"context"

	"github.com/humio/cli/internal/api/humiographql"
)

type Users struct {
	client *Client
}

type User struct {
	ID          string
	Username    string
	FullName    *string
	Email       *string
	Company     *string
	CountryCode *string
	Picture     *string
	IsRoot      bool
	CreatedAt   string
}

func (c *Client) Users() *Users { return &Users{client: c} }

func (u *Users) List() ([]User, error) {
	resp, err := humiographql.ListUsers(context.Background(), u.client)
	if err != nil {
		return nil, err
	}

	respUsers := resp.GetUsers()
	users := make([]User, len(respUsers))
	for idx, user := range respUsers {
		users[idx] = User{
			ID:          user.GetId(),
			Username:    user.GetUsername(),
			FullName:    user.GetFullName(),
			Email:       user.GetEmail(),
			Company:     user.GetCompany(),
			CountryCode: user.GetCountryCode(),
			Picture:     user.GetPicture(),
			IsRoot:      user.GetIsRoot(),
			CreatedAt:   user.GetCreatedAt().String(),
		}
	}

	return users, nil
}

func (u *Users) Get(username string) (User, error) {
	resp, err := humiographql.GetUsersByUsername(context.Background(), u.client, username)
	if err != nil {
		return User{}, err
	}

	respUsers := resp.GetUsers()
	for _, user := range respUsers {
		if user.Username == username {
			return User{
				ID:          user.GetId(),
				Username:    user.GetUsername(),
				FullName:    user.GetFullName(),
				Email:       user.GetEmail(),
				Company:     user.GetCompany(),
				CountryCode: user.GetCountryCode(),
				Picture:     user.GetPicture(),
				IsRoot:      user.GetIsRoot(),
				CreatedAt:   user.GetCreatedAt().String(),
			}, nil
		}
	}

	return User{}, UserNotFound(username)
}

func (u *Users) Update(username string, isRoot *bool, fullName, company, countryCode, email, picture *string) (User, error) {
	_, err := humiographql.UpdateUser(context.Background(), u.client, username, company, isRoot, fullName, picture, email, countryCode)
	if err != nil {
		return User{}, err
	}

	return u.Get(username)

}

func (u *Users) Add(username string, isRoot *bool, fullName, company, countryCode, email, picture *string) (User, error) {
	resp, err := humiographql.AddUser(context.Background(), u.client, username, company, isRoot, fullName, picture, email, countryCode)
	if err != nil {
		return User{}, err
	}

	createdUser := resp.GetAddUserV2()
	switch v := createdUser.(type) {
	case *humiographql.AddUserAddUserV2User:
		return User{
			ID:          v.GetId(),
			Username:    v.GetUsername(),
			FullName:    v.GetFullName(),
			Email:       v.GetEmail(),
			Company:     v.GetCompany(),
			CountryCode: v.GetCountryCode(),
			Picture:     v.GetPicture(),
			IsRoot:      v.GetIsRoot(),
			CreatedAt:   v.GetCreatedAt().String(),
		}, nil
	default:
		panic("not implemented")
	}
}

func (u *Users) Remove(username string) (User, error) {
	resp, err := humiographql.RemoveUser(context.Background(), u.client, username)
	if err != nil {
		return User{}, err
	}
	respUser := resp.GetRemoveUser()
	user := respUser.GetUser()
	return User{
		ID:          user.GetId(),
		Username:    user.GetUsername(),
		FullName:    user.GetFullName(),
		Email:       user.GetEmail(),
		Company:     user.GetCompany(),
		CountryCode: user.GetCountryCode(),
		Picture:     user.GetPicture(),
		IsRoot:      user.GetIsRoot(),
		CreatedAt:   user.GetCreatedAt().String(),
	}, nil
}
