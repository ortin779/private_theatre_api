package repository

import (
	"database/sql"
	"fmt"

	"github.com/ortin779/private_theatre_api/api/models"
)

type AddonRepository interface {
	Create(addon models.Addon) error
	GetCategories() []string
	GetAllAddons() ([]models.Addon, error)
}

type addonRepository struct {
	db *sql.DB
}

func NewAddonRepository(db *sql.DB) AddonRepository {
	return &addonRepository{
		db: db,
	}
}

func (as *addonRepository) Create(addon models.Addon) error {
	_, err := as.db.Exec(`INSERT INTO addons(id, name, category, price, meta_data, created_at, updated_at, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, addon.ID, addon.Name, addon.Category, addon.Price, addon.MetaData, addon.CreatedAt, addon.UpdatedAt, addon.CreatedBy, addon.UpdatedBy)

	if err != nil {
		return fmt.Errorf("create addon: %w", err)
	}
	return nil
}

func (as *addonRepository) GetCategories() []string {
	return models.AddonCategories
}

func (as *addonRepository) GetAllAddons() ([]models.Addon, error) {
	rows, err := as.db.Query(`SELECT * FROM addons;`)
	if err != nil {
		return nil, fmt.Errorf("get addons: %w", err)
	}
	defer rows.Close()
	var addons []models.Addon

	for rows.Next() {
		var addon models.Addon
		err = rows.Scan(&addon.ID, &addon.Name, &addon.Category, &addon.Price, &addon.MetaData, &addon.UpdatedAt, &addon.CreatedAt, &addon.CreatedBy, &addon.UpdatedBy)
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
