package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ortin779/private_theatre_api/api/models"
)

var (
	ErrNoUserWithEmail = errors.New("no user found with given email id")
	ErrNoUserWithId    = errors.New("no user found with given user id")
)

type UsersRepository interface {
	Create(user models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByUserId(id string) (*models.User, error)
}

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (ur *usersRepository) Create(user models.User) error {
	_, err := ur.db.Exec(`INSERT INTO users(id, name, email, password, roles)
    VALUES($1,$2,$3,$4,$5);
`, user.ID, user.Name, user.Email, user.Password, user.Roles)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (ur *usersRepository) GetByEmail(email string) (*models.User, error) {
	row := ur.db.QueryRow(`SELECT * FROM users
		WHERE email=$1;`, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, pgtype.NewMap().SQLScanner(&user.Roles))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoUserWithEmail
		}
		return nil, err
	}
	return &user, nil
}

func (ur *usersRepository) GetByUserId(id string) (*models.User, error) {
	row := ur.db.QueryRow(`SELECT * FROM users
		WHERE id=$1;`, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Roles)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoUserWithId
		}
		return nil, err
	}

	return &user, nil
}
