package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	MinTime = 0
	MaxTime = 1440
)

type Slot struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	CreatedBy string    `json:"created_by"`
}

type CreateSlotParams struct {
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
}

func (csp CreateSlotParams) Validate() map[string]string {
	errs := map[string]string{}
	if csp.StartTime < MinTime {
		errs["start_time"] = fmt.Sprintf("start time should be minimum of %d", MinTime)
	}
	if csp.StartTime > MaxTime {
		errs["start_time"] = fmt.Sprintf("start time should be maximum of %d", MaxTime)
	}
	if csp.EndTime < MinTime {
		errs["start_time"] = fmt.Sprintf("start time should be minimum of %d", MinTime)
	}
	if csp.EndTime > MaxTime {
		errs["start_time"] = fmt.Sprintf("start time should be maximum of %d", MaxTime)
	}
	if csp.StartTime >= csp.EndTime {
		errs["start_time"] = fmt.Sprintf("start time: %d, should be lessthan endtime: %d", csp.StartTime, csp.EndTime)
	}

	return errs
}

type SlotStore interface {
	GetSlots() ([]Slot, error)
	AddSlot(slot Slot) error
}

func NewSlotStore(db *sql.DB) SlotStore {
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
		err := slotRows.Scan(&slot.ID, &slot.StartTime, &slot.EndTime, &slot.UpdatedAt, &slot.CreatedAt, &slot.CreatedBy, &slot.UpdatedBy)
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
	_, err := ss.db.Exec(`
		INSERT INTO slots(id, start_time, end_time, created_at, updated_at, created_by, updated_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, slot.ID, slot.StartTime, slot.EndTime, slot.CreatedAt, slot.UpdatedAt, slot.CreatedBy, slot.UpdatedBy)
	// TODO: Handle custom errors
	return err
}
