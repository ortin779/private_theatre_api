package models

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type CreateTheatreParams struct {
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	Price                  float64  `json:"price"`
	AdditionalPricePerHead float64  `json:"additional_price_per_head"`
	MaxCapacity            int      `json:"max_capacity"`
	MinCapacity            int      `json:"min_capacity"`
	DefaultCapacity        int      `json:"default_capacity"`
	Slots                  []string `json:"slots"`
}

func (ctp CreateTheatreParams) Validate() map[string]string {
	errors := make(map[string]string)

	if ctp.Name == "" {
		errors["name"] = "name of the theatre can not be empty"
	}

	if ctp.Description == "" {
		errors["description"] = "description of the theatre can not be empty"
	}

	if ctp.Price <= 0 {
		errors["price"] = "price of the theatre can not be zero or negative"
	}

	if ctp.AdditionalPricePerHead <= 0 {
		errors["additional_price_per_head"] = "additional price per head should be a positive number"
	}

	if len(ctp.Slots) == 0 {
		errors["slots"] = "theatre should have at least one slot allocated"
	}

	for _, val := range ctp.Slots {
		if _, err := uuid.Parse(val); err != nil {
			errors["slots"] = "invalid slot_id, it should be a valid uuid"
			break
		}
	}
	return errors
}

type Theatre struct {
	ID                     string  `json:"id"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	Price                  float64 `json:"price"`
	AdditionalPricePerHead float64 `json:"additional_price_per_head"`
	MaxCapacity            int     `json:"max_capacity"`
	MinCapacity            int     `json:"min_capacity"`
	DefaultCapacity        int     `json:"default_capacity"`
}

type TheatreStore interface {
	GetTheatres() ([]Theatre, error)
	Create(t Theatre, slots []string) error
}

type TheatreService struct {
	db *sql.DB
}

func NewTheatreService(db *sql.DB) *TheatreService {
	return &TheatreService{
		db: db,
	}
}

func (ts *TheatreService) GetTheatres() ([]Theatre, error) {
	return nil, nil
}

func (ts *TheatreService) Create(t Theatre, slots []string) error {
	tx, err := ts.db.Begin()

	if err != nil {
		return fmt.Errorf("create theatre: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        INSERT INTO theatres(id, name, description, price, additional_price_per_head, max_capacity, min_capacity, default_capacity) Values ($1, $2, $3, $4, $5, $6, $7,$8);
    `, t.ID, t.Name, t.Description, t.Price, t.AdditionalPricePerHead, t.MaxCapacity, t.MinCapacity, t.DefaultCapacity)

	if err != nil {
		return fmt.Errorf("create theatre: %w", err)
	}

	stmt, err := tx.Prepare(`
        INSERT INTO theatre_slots(theatre_id, slot_id) VALUES ($1, $2);
    `)
	if err != nil {
		return fmt.Errorf("create theatre: %w", err)
	}

	for _, slotId := range slots {
		_, err := stmt.Exec(t.ID, slotId)
		if err != nil {
			return fmt.Errorf("create theatre: %w", err)
		}
	}

	err = tx.Commit()

	return err
}
