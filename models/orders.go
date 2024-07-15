package models

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type OrderAddon struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderParams struct {
	CustomerName  string       `json:"customer_name"`
	CustomerEmail string       `json:"customer_email"`
	PhoneNumber   string       `json:"phone_number"`
	TheatreId     string       `json:"theatre_id"`
	SlotId        string       `json:"slot_id"`
	NoOfPersons   int          `json:"no_of_persons"`
	TotalPrice    float64      `json:"total_price"`
	OrderDate     time.Time    `json:"order_date"`
	Addons        []OrderAddon `json:"addons"`
}

func (op *OrderParams) Validate() map[string]string {
	errs := make(map[string]string)

	if len(op.PhoneNumber) == 0 {
		errs["phone_number"] = "phone number can not be empty"
	}
	if len(op.CustomerName) == 0 {
		errs["customer_name"] = "phone number can not be empty"
	}
	if !isEmailValid(op.CustomerEmail) {
		errs["customer_email"] = "invalid email address"
	}
	if _, err := uuid.Parse(op.SlotId); err != nil {
		errs["slot_id"] = "slot id must be a valid uuid"
	}
	if op.TotalPrice <= 0 {
		errs["total_price"] = "order value must be greater than zero"
	}
	for _, addon := range op.Addons {
		if _, err := uuid.Parse(addon.ID); err != nil {
			errs["addons"] = "addon id must be a valid uuid"
			break
		}
		if addon.Quantity <= 0 {
			errs["addons"] = "addon quantity should not be less than 1"
			break
		}
	}
	return errs
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type OrderDetails struct {
	ID            string    `json:"id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	PhoneNumber   string    `json:"phone_number"`
	NoOfPersons   int       `json:"no_of_persons"`
	TotalPrice    float64   `json:"total_price"`
	OrderDate     time.Time `json:"order_date"`
	Theatre       Theatre   `json:"theatre"`
	Slot          Slot      `json:"slot"`
	Addons        []Addon   `json:"addons"`
	OrderedAt     time.Time `json:"ordered_at"`
}

type Order struct {
	ID            string    `json:"id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	PhoneNumber   string    `json:"phone_number"`
	TheatreId     string    `json:"theatre_id"`
	SlotId        string    `json:"slot_id"`
	NoOfPersons   int       `json:"no_of_persons"`
	TotalPrice    float64   `json:"total_price"`
	OrderDate     time.Time `json:"order_date"`
	OrderedAt     time.Time `json:"ordered_at"`
}

type OrderStore interface {
	Create(orderParams OrderParams) (Order, error)
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

var ErrDuplicateOrder = errors.New("theatre already booked for that day and slot")

func (orderService *OrderService) Create(orderParams OrderParams) (Order, error) {
	tx, err := orderService.db.Begin()

	if err != nil {
		return Order{}, fmt.Errorf("create order: %w", err)
	}
	defer tx.Rollback()

	order := Order{
		ID:            uuid.NewString(),
		CustomerName:  orderParams.CustomerName,
		CustomerEmail: orderParams.CustomerEmail,
		PhoneNumber:   orderParams.PhoneNumber,
		TotalPrice:    orderParams.TotalPrice,
		OrderDate:     orderParams.OrderDate,
		NoOfPersons:   orderParams.NoOfPersons,
		SlotId:        orderParams.SlotId,
		TheatreId:     orderParams.TheatreId,
	}

	row := tx.QueryRow(`SELECT id FROM orders
        WHERE theatre_id = $1 AND slot_id = $2 AND order_date = $3;
    `, order.TheatreId, order.SlotId, order.OrderDate.Format(time.DateOnly))

	err = row.Scan()
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return Order{}, ErrDuplicateOrder
	}

	row = tx.QueryRow(`INSERT INTO orders(
    id,customer_name,customer_email,phone_number,no_of_persons,total_price,order_date,theatre_id, slot_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING ordered_at;`, order.ID, order.CustomerName, order.CustomerEmail, order.PhoneNumber, order.NoOfPersons, order.TotalPrice, order.OrderDate.Format(time.DateOnly), order.TheatreId, order.SlotId)

	if err := row.Scan(&order.OrderedAt); err != nil {
		return Order{}, fmt.Errorf("create order: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO order_addons(order_id, addon_id, quantity) VALUES ($1,$2,$3);")
	if err != nil {
		return Order{}, fmt.Errorf("create order: %w", err)
	}

	for _, addon := range orderParams.Addons {
		_, err = stmt.Exec(order.ID, addon.ID, addon.Quantity)
		if err != nil {
			return Order{}, fmt.Errorf("create order: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return Order{}, fmt.Errorf("create order: %w", err)
	}

	return order, nil
}

func (orderService *OrderService) GetAll() ([]OrderDetails, error) {
	return []OrderDetails{}, nil
}

func (orderSerice *OrderService) GetById(id string) (OrderDetails, error) {
	return OrderDetails{}, nil
}
