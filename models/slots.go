package models

import (
	"database/sql"
	"log"
)

type Slot struct {
	ID        string
	StartTime int
	EndTime   int
}

type SlotStore interface {
	GetSlots() ([]Slot, error)
	AddSlot(slot Slot) error
	UpdateSlot(id int, slot Slot) error
	DeleteSlot(id int) error
}

func NewSlotService(db *sql.DB) *SlotsService {
	return &SlotsService{
		db: db,
	}
}

type SlotsService struct {
	db *sql.DB
}

func (ss *SlotsService) GetSlots() ([]Slot, error) {
	var slots []Slot
	slotRows, err := ss.db.Query(`SELECT * FROM slots;`)
	if err != nil {
		return nil, err
	}
	defer slotRows.Close()

	for slotRows.Next() {
		var slot Slot
		err := slotRows.Scan(&slot.ID, slot.StartTime, slot.EndTime)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}

	err = slotRows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return slots, nil
}

func (ss *SlotsService) AddSlot(slot Slot) error {
	return nil
}

func (ss *SlotsService) UpdateSlot(id int, slot Slot) error {
	return nil
}

func (ss *SlotsService) DeleteSlot(id int) error {
	return nil
}
