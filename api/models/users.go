package models

import (
	"errors"
	"slices"
)

var UserRoles = []string{
	"admin",
	// this role is for future purpose, where we need to have authentication
	// for the normal users as well
	"customer",
}

var (
	ErrNoUserWithEmail = errors.New("no user found with given email id")
	ErrNoUserWithId    = errors.New("no user found with given user id")
)

type UserParams struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

func (userParams UserParams) Validate() map[string]string {
	errs := make(map[string]string)

	if userParams.Name == "" {
		errs["name"] = "user name can not be empty"
	}
	if !isEmailValid(userParams.Email) {
		errs["email"] = "user email should be a valid email address"
	}
	if len(userParams.Password) < 8 {
		errs["password"] = "user password should be at least 8 characters"
	}
	for _, role := range userParams.Roles {
		if !slices.Contains(UserRoles, role) {
			errs["roles"] = role + " is not a valid role"
			break
		}
	}
	return errs
}

type User struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"-"`
	Roles    []string `json:"roles"`
}
