package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"slices"
	"time"
)

type AddonCategory int

const (
	Decorations AddonCategory = iota + 1
	Cakes
	Flowers
	Photographs
)

var AddonCategories = []string{Decorations.String(), Flowers.String(), Cakes.String(), Photographs.String()}

func (ac AddonCategory) String() string {
	switch ac {
	case 1:
		return "Decorations"
	case 2:
		return "Cakes"
	case 3:
		return "Flowers"
	case 4:
		return "Photographs"
	default:
		return ""
	}
}

type MetaData map[string]any

func (m MetaData) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MetaData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &m)
}

type AddonParams struct {
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Price    float64  `json:"price"`
	MetaData MetaData `json:"meta_data,omitempty"`
}

func (ap *AddonParams) Validate() map[string]string {
	errors := make(map[string]string)
	if ap.Name == "" {
		errors["name"] = "addon name can not be empty"
	}
	if ap.Price <= 0 {
		errors["price"] = "addon price can not be negative"
	}
	if !slices.Contains(AddonCategories, ap.Category) {
		errors["category"] = "addon category is not valid"
	}
	return errors
}

type Addon struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Price     float64   `json:"price"`
	MetaData  MetaData  `json:"meta_data,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	CreatedBy string    `json:"created_by"`
}
