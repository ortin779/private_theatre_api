package models

import "database/sql"

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
	Create(t Theatre) error
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

func (ts *TheatreService) Create() error {
	return nil
}
