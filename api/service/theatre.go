package service

import (
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/repository"
)

type TheatresService struct {
	theatresRepo repository.TheatreRepository
}

func NewTheatreService(theatresRepo repository.TheatreRepository) TheatresService {
	return TheatresService{
		theatresRepo: theatresRepo,
	}
}

func (ts *TheatresService) Create(t models.Theatre, slots []string) error {
	return ts.theatresRepo.Create(t, slots)
}

func (ts *TheatresService) GetTheatres() ([]models.Theatre, error) {
	return ts.theatresRepo.GetTheatres()
}

func (ts *TheatresService) GetTheatreDetails(id string) (*models.TheatreWithSlots, error) {
	return ts.theatresRepo.GetTheatreDetails(id)
}
