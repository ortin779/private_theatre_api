package repository

import (
	"database/sql"
	"fmt"

	"github.com/ortin779/private_theatre_api/api/models"
)

type TheatreRepository interface {
	GetTheatres() ([]models.Theatre, error)
	Create(t models.Theatre, slots []string) error
	GetTheatreDetails(id string) (*models.TheatreWithSlots, error)
}

type theatreRepository struct {
	db *sql.DB
}

func NewTheatreRepository(db *sql.DB) TheatreRepository {
	return &theatreRepository{
		db: db,
	}
}

func (tr *theatreRepository) GetTheatres() ([]models.Theatre, error) {
	var theatres []models.Theatre
	rows, err := tr.db.Query(`
		SELECT * FROM theatres;
	`)
	if err != nil {
		return nil, fmt.Errorf("get theatres: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var theatre models.Theatre
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

func (tr *theatreRepository) Create(t models.Theatre, slots []string) error {
	tx, err := tr.db.Begin()

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

func (tr *theatreRepository) GetTheatreDetails(id string) (*models.TheatreWithSlots, error) {
	var theatreDetails models.TheatreWithSlots
	row := tr.db.QueryRow(`
		SELECT * FROM theatres
			WHERE id = $1;
	`, id)

	err := row.Scan(&theatreDetails.ID, &theatreDetails.Name, &theatreDetails.Description, &theatreDetails.Price, &theatreDetails.AdditionalPricePerHead, &theatreDetails.MaxCapacity, &theatreDetails.MinCapacity, &theatreDetails.DefaultCapacity, &theatreDetails.UpdatedAt, &theatreDetails.CreatedAt, &theatreDetails.CreatedBy, &theatreDetails.UpdatedBy)

	if err != nil {
		return nil, fmt.Errorf("get theatre details: %w", err)
	}

	var slots []models.Slot

	rows, err := tr.db.Query(`
		SELECT * FROM slots
			WHERE id IN (
			SELECT slot_id from theatre_slots WHERE theatre_id=$1
			);
	`, id)

	if err != nil {
		return nil, fmt.Errorf("get theatre details: %w", err)
	}

	for rows.Next() {
		var slot models.Slot
		err := rows.Scan(&slot.ID, &slot.StartTime, &slot.EndTime, &slot.UpdatedAt, &slot.CreatedAt, &slot.CreatedBy, &slot.UpdatedBy)
		if err != nil {
			return nil, fmt.Errorf("get theatre details: %w", err)
		}
		slots = append(slots, slot)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("get theatre details: %w", row.Err())
	}
	theatreDetails.Slots = slots

	return &theatreDetails, nil
}
