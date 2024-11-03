package models

import (
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
