package service

import (
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/repository"
)

type UsersService struct {
	usersRepo repository.UsersRepository
}

func NewUsersService(usersRepo repository.UsersRepository) UsersService {
	return UsersService{
		usersRepo: usersRepo,
	}
}

func (us *UsersService) Create(user models.User) error {
	return us.usersRepo.Create(user)
}

func (us *UsersService) GetByEmail(email string) (*models.User, error) {
	return us.usersRepo.GetByEmail(email)
}

func (us *UsersService) GetByUserId(userId string) (*models.User, error) {
	return us.usersRepo.GetByUserId(userId)
}
