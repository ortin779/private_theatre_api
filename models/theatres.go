package models

import (
	"database/sql"
	"fmt"
	"time"

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
	ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	Price                  float64   `json:"price"`
	AdditionalPricePerHead float64   `json:"additional_price_per_head"`
	MaxCapacity            int       `json:"max_capacity"`
	MinCapacity            int       `json:"min_capacity"`
	DefaultCapacity        int       `json:"default_capacity"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	UpdatedBy              string    `json:"updated_by"`
	CreatedBy              string    `json:"created_by"`
}

type TheatreWithSlots struct {
	Theatre
	Slots []Slot `json:"slots"`
}

type TheatreStore interface {
	GetTheatres() ([]Theatre, error)
	Create(t Theatre, slots []string) error
	GetTheatreDetails(id string) (TheatreWithSlots, error)
}

type TheatreService struct {
	db *sql.DB
}

func NewTheatreStore(db *sql.DB) TheatreStore {
	return &TheatreService{
		db: db,
	}
}

func (ts *TheatreService) GetTheatres() ([]Theatre, error) {
	var theatres []Theatre
	rows, err := ts.db.Query(`
		SELECT * FROM theatres;
	`)
	if err != nil {
		return nil, fmt.Errorf("get theatres: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var theatre Theatre
		err := rows.Scan(&theatre.ID, &theatre.Name, &theatre.Description, &theatre.Price, &theatre.AdditionalPricePerHead, &theatre.MaxCapacity, &theatre.MinCapacity, &theatre.DefaultCapacity, &theatre.UpdatedAt, &theatre.CreatedAt, &theatre.CreatedBy, &theatre.UpdatedBy)
		if err != nil {
			return nil, fmt.Errorf("get theatres: %w", err)
		}
		theatres = append(theatres, theatre)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("get theatres, %w", rows.Err())
	}
	return theatres, nil
}

func (ts *TheatreService) Create(t Theatre, slots []string) error {
	tx, err := ts.db.Begin()

	if err != nil {
		return fmt.Errorf("create theatre: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        INSERT INTO theatres(id, name, description, price, additional_price_per_head, max_capacity, min_capacity, default_capacity, created_at, updated_at, created_by, updated_by) Values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
    `, t.ID, t.Name, t.Description, t.Price, t.AdditionalPricePerHead, t.MaxCapacity, t.MinCapacity, t.DefaultCapacity, t.CreatedAt, t.UpdatedAt, t.CreatedBy, t.UpdatedBy)

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

func (ts *TheatreService) GetTheatreDetails(id string) (TheatreWithSlots, error) {
	var theatreDetails TheatreWithSlots
	row := ts.db.QueryRow(`
		SELECT * FROM theatres
			WHERE id = $1;
	`, id)

	err := row.Scan(&theatreDetails.ID, &theatreDetails.Name, &theatreDetails.Description, &theatreDetails.Price, &theatreDetails.AdditionalPricePerHead, &theatreDetails.MaxCapacity, &theatreDetails.MinCapacity, &theatreDetails.DefaultCapacity, &theatreDetails.UpdatedAt, &theatreDetails.CreatedAt, &theatreDetails.CreatedBy, &theatreDetails.UpdatedBy)

	if err != nil {
		return TheatreWithSlots{}, fmt.Errorf("get theatre details: %w", err)
	}

	var slots []Slot

	rows, err := ts.db.Query(`
		SELECT * FROM slots
			WHERE id IN (
			SELECT slot_id from theatre_slots WHERE theatre_id=$1
			);
	`, id)

	if err != nil {
		return TheatreWithSlots{}, fmt.Errorf("get theatre details: %w", err)
	}

	for rows.Next() {
		var slot Slot
		err := rows.Scan(&slot.ID, &slot.StartTime, &slot.EndTime, &slot.UpdatedAt, &slot.CreatedAt, &slot.CreatedBy, &slot.UpdatedBy)
		if err != nil {
			return TheatreWithSlots{}, fmt.Errorf("get theatre details: %w", err)
		}
		slots = append(slots, slot)
	}

	if rows.Err() != nil {
		return TheatreWithSlots{}, fmt.Errorf("get theatre details: %w", row.Err())
	}
	theatreDetails.Slots = slots

	return theatreDetails, nil
}
