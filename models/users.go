package models

import (
	"database/sql"
	"fmt"
	"slices"
)

var UserRoles = []string{
	"admin",
	// this role is for future purpose, where we need to have authentication
	// for the normal users as well
	"customer",
}

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

type UserStore interface {
	Create(user User) error
	GetByEmail(email string) (User, error)
}

type UserService struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return &UserService{
		db: db,
	}
}

func (usrService *UserService) Create(user User) error {
	_, err := usrService.db.Exec(`INSERT INTO users(id, name, email, password, roles)
        VALUES($1,$2,$3,$4,$5);
    `, user.ID, user.Name, user.Email, user.Password, user.Roles)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (usrService *UserService) GetByEmail(email string) (User, error) {
	return User{}, nil
}
