package repository

import (
	"database/sql"

	"github.com/ortin779/private_theatre_api/api/models"
)

type SlotsRepository interface {
	GetSlots() ([]models.Slot, error)
	AddSlot(slot models.Slot) error
}

type slotsRepository struct {
	db *sql.DB
}

func NewSlotsRepo(db *sql.DB) SlotsRepository {
	return &slotsRepository{
		db: db,
	}
}

func (sr *slotsRepository) GetSlots() ([]models.Slot, error) {
	var slots []models.Slot
	slotRows, err := sr.db.Query(`SELECT * FROM slots;`)
	if err != nil {
		return nil, err
	}
	defer slotRows.Close()

	for slotRows.Next() {
		var slot models.Slot
		err := slotRows.Scan(&slot.ID, &slot.StartTime, &slot.EndTime, &slot.UpdatedAt, &slot.CreatedAt, &slot.CreatedBy, &slot.UpdatedBy)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	err = slotRows.Err()
	if err != nil {
		return nil, err
	}

	return slots, nil
}

func (sr *slotsRepository) AddSlot(slot models.Slot) error {
	_, err := sr.db.Exec(`
		INSERT INTO slots(id, start_time, end_time, created_at, updated_at, created_by, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, slot.ID, slot.StartTime, slot.EndTime, slot.CreatedAt, slot.UpdatedAt, slot.CreatedBy, slot.UpdatedBy)
	// TODO: Handle custom errors
	return err
}
