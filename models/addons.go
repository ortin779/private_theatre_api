package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
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
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Price    float64  `json:"price"`
	MetaData MetaData `json:"meta_data,omitempty"`
}

type AddonStore interface {
	Create(addon Addon) error
	GetCategories() []string
	GetAddons() ([]Addon, error)
}

type AddonService struct {
	db *sql.DB
}

func NewAddonStore(db *sql.DB) *AddonService {
	return &AddonService{
		db: db,
	}
}

func (as *AddonService) Create(addon Addon) error {
	_, err := as.db.Exec(`INSERT INTO addons(id, name, category, price, meta_data)
        VALUES ($1, $2, $3, $4, $5)
    `, addon.ID, addon.Name, addon.Category, addon.Price, addon.MetaData)

	if err != nil {
		return fmt.Errorf("create addon: %w", err)
	}
	return nil
}

func (as *AddonService) GetCategories() []string {
	return AddonCategories
}

func (as *AddonService) GetAddons() ([]Addon, error) {
	rows, err := as.db.Query(`SELECT * FROM addons;`)
	if err != nil {
		return nil, fmt.Errorf("get addons: %w", err)
	}
	defer rows.Close()
	var addons []Addon

	for rows.Next() {
		var addon Addon
		err = rows.Scan(&addon.ID, &addon.Name, &addon.Category, &addon.Price, &addon.MetaData)
		if err != nil {
			return nil, fmt.Errorf("get addons: %w", err)
		}
		addons = append(addons, addon)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("get addons: %w", rows.Err())
	}

	return addons, nil
}
