package service

import (
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/repository"
)

type AddonsService struct {
	addonsRepo repository.AddonRepository
}

func NewAddonService(addonsRepo repository.AddonRepository) AddonsService {
	return AddonsService{
		addonsRepo: addonsRepo,
	}
}

func (as *AddonsService) CreateAddon(addon models.Addon) error {
	return as.addonsRepo.Create(addon)
}

func (as *AddonsService) GetCategories() []string {
	return as.addonsRepo.GetCategories()
}

func (as *AddonsService) GetAllAddons() ([]models.Addon, error) {
	return as.addonsRepo.GetAllAddons()
}
