package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type OrderParams struct {
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	PhoneNumber   string    `json:"phone_number"`
	TheatreId     string    `json:"theatre_id"`
	SlotId        string    `json:"slot_id"`
	NoOfPersons   int       `json:"no_of_persons"`
	TotalPrice    float64   `json:"total_price"`
	OrderDate     time.Time `json:"order_date"`
	Addons        []string  `json:"addons"`
}

func (op *OrderParams) Validate() map[string]string {
	errs := make(map[string]string)

	if len(op.PhoneNumber) == 0 {
		errs["phone_number"] = "phone number can not be empty"
	}
	if _, err := uuid.Parse(op.TheatreId); err != nil {
		errs["theatre_id"] = "theatre id should be a valid uuid"
	}
	if _, err := uuid.Parse(op.SlotId); err != nil {
		errs["slot_id"] = "slot id must be a valid uuid"
	}
	if op.TotalPrice <= 0 {
		errs["total_price"] = "order value must be greater than zero"
	}
	for _, addonId := range op.Addons {
		if _, err := uuid.Parse(addonId); err != nil {
			errs["addons"] = "addon id must be a valid uuid"
			break
		}
	}
	return errs
}

type OrderDetails struct {
	ID            string    `json:"id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	PhoneNumber   string    `json:"phone_number"`
	NoOfPersons   int       `json:"no_of_persons"`
	Total_Price   float64   `json:"total_price"`
	OrderDate     time.Time `json:"order_date"`
	Theatre       Theatre   `json:"theatre"`
	Slot          Slot      `json:"slot"`
	Addons        []Addon   `json:"addons"`
}

type Order struct {
	ID            string    `json:"id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	PhoneNumber   string    `json:"phone_number"`
	TheatreId     string    `json:"theatre_id"`
	SlotId        string    `json:"slot_id"`
	NoOfPersons   int       `json:"no_of_persons"`
	Total_Price   float64   `json:"total_price"`
	OrderDate     time.Time `json:"order_date"`
}

type OrderStore interface {
	Create(op OrderParams) error
	GetAll() ([]OrderDetails, error)
	GetById(id string) (OrderDetails, error)
}

type OrderService struct {
	db *sql.DB
}

func NewOrderStore(db *sql.DB) OrderStore {
	return &OrderService{
		db: db,
	}
}

func (orderService *OrderService) Create(orderParams OrderParams) error {
	return nil
}

func (orderService *OrderService) GetAll() ([]OrderDetails, error) {
	return []OrderDetails{}, nil
}

func (orderSerice *OrderService) GetById(id string) (OrderDetails, error) {
	return OrderDetails{}, nil
}
