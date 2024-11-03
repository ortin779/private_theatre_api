package service

import (
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/repository"
)

type SlotsService struct {
	slotsRepo repository.SlotsRepository
}

func NewSlotsService(slotsRepo repository.SlotsRepository) SlotsService {
	return SlotsService{
		slotsRepo: slotsRepo,
	}
}

func (ss *SlotsService) GetSlots() ([]models.Slot, error) {
	return ss.slotsRepo.GetSlots()
}

func (ss *SlotsService) AddSlot(slot models.Slot) error {
	return ss.slotsRepo.AddSlot(slot)
}
