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

type OrderAddonDetails struct {
	Addon
	Quantity int `json:"quantity"`
}

type PaymentStatus string

var (
	Success PaymentStatus = "success"
	Failure PaymentStatus = "failure"
	Pending PaymentStatus = "pending"
)

type OrderParams struct {
	CustomerName  string       `json:"customer_name"`
	CustomerEmail string       `json:"customer_email"`
	PhoneNumber   string       `json:"phone_number"`
	TheatreId     string       `json:"theatre_id"`
	SlotId        string       `json:"slot_id"`
	NoOfPersons   int          `json:"no_of_persons"`
	TotalPrice    int          `json:"total_price"`
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
	ID             string              `json:"id"`
	CustomerName   string              `json:"customer_name"`
	CustomerEmail  string              `json:"customer_email"`
	PhoneNumber    string              `json:"phone_number"`
	NoOfPersons    int                 `json:"no_of_persons"`
	TotalPrice     int                 `json:"total_price"`
	OrderDate      time.Time           `json:"order_date"`
	Theatre        Theatre             `json:"theatre"`
	Slot           Slot                `json:"slot"`
	Addons         []OrderAddonDetails `json:"addons"`
	OrderedAt      time.Time           `json:"ordered_at"`
	PaymentDetails OrderPayment        `json:"payment_details"`
}

type Order struct {
	ID              string       `json:"id"`
	CustomerName    string       `json:"customer_name"`
	CustomerEmail   string       `json:"customer_email"`
	PhoneNumber     string       `json:"phone_number"`
	TheatreId       string       `json:"theatre_id"`
	SlotId          string       `json:"slot_id"`
	Addons          []OrderAddon `json:"addons"`
	NoOfPersons     int          `json:"no_of_persons"`
	TotalPrice      int          `json:"total_price"`
	OrderDate       time.Time    `json:"order_date"`
	OrderedAt       time.Time    `json:"ordered_at"`
	RazorpayOrderId string       `json:"razorpay_order_id"`
}

type OrderStore interface {
	Create(order Order) error
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

func (orderService *OrderService) Create(order Order) error {
	tx, err := orderService.db.Begin()

	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT id FROM orders
        WHERE theatre_id = $1 AND slot_id = $2 AND order_date = $3;
    `, order.TheatreId, order.SlotId, order.OrderDate.Format(time.DateOnly))

	err = row.Scan()
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		return ErrDuplicateOrder
	}

	row = tx.QueryRow(`INSERT INTO orders(
    id,customer_name,customer_email,phone_number,no_of_persons,total_price,order_date,theatre_id, slot_id, razorpay_order_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING ordered_at;`, order.ID, order.CustomerName, order.CustomerEmail, order.PhoneNumber, order.NoOfPersons, order.TotalPrice, order.OrderDate.Format(time.DateOnly), order.TheatreId, order.SlotId, order.RazorpayOrderId)

	if err := row.Scan(&order.OrderedAt); err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO order_addons(order_id, addon_id, quantity) VALUES ($1,$2,$3);")
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	for _, addon := range order.Addons {
		_, err = stmt.Exec(order.ID, addon.ID, addon.Quantity)
		if err != nil {
			return fmt.Errorf("create order: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return nil
}

func (orderService *OrderService) GetAll() ([]OrderDetails, error) {

	rows, err := orderService.db.Query(`SELECT
		orders.id,
		orders.customer_name,
		orders.customer_email,
		orders.phone_number,
		orders.no_of_persons,
		orders.total_price,
		orders.order_date,
		orders.ordered_at,
		theatres.id,
		theatres."name" ,
		theatres.description ,
		theatres.price ,
		theatres.additional_price_per_head ,
		theatres.max_capacity ,
		theatres.min_capacity ,
		theatres.default_capacity ,
		slots.id ,
		slots.start_time ,
		slots.end_time,
		payments.razorpay_order_id,
		payments.razorpay_payment_id,
		payments.razorpay_signature,
		payments.status
	FROM
		orders
	JOIN theatres ON
		orders.theatre_id = theatres.id
	JOIN slots ON
		slots.id = orders.slot_id
	JOIN payments ON
		orders.razorpay_order_id = payments.razorpay_order_id;`)

	if err != nil {
		return nil, fmt.Errorf("get orders: %w", err)
	}
	defer rows.Close()

	orderDetailsList := make([]OrderDetails, 0, 5)
	for rows.Next() {
		var orderDetails OrderDetails
		err := rows.Scan(&orderDetails.ID, &orderDetails.CustomerName, &orderDetails.CustomerEmail, &orderDetails.PhoneNumber, &orderDetails.NoOfPersons, &orderDetails.TotalPrice, &orderDetails.OrderDate, &orderDetails.OrderedAt, &orderDetails.Theatre.ID, &orderDetails.Theatre.Name, &orderDetails.Theatre.Description, &orderDetails.Theatre.Price, &orderDetails.Theatre.AdditionalPricePerHead, &orderDetails.Theatre.MaxCapacity, &orderDetails.Theatre.MinCapacity, &orderDetails.Theatre.DefaultCapacity, &orderDetails.Slot.ID, &orderDetails.Slot.StartTime, &orderDetails.Slot.EndTime, &orderDetails.PaymentDetails.RazorpayOrderId, &orderDetails.PaymentDetails.RazorpayPaymentId, &orderDetails.PaymentDetails.RazorpaySignature, &orderDetails.PaymentDetails.Status)

		if err != nil {
			return nil, fmt.Errorf("get orders: %w", err)
		}

		addons, err := orderService.getAddonsForOrder(orderDetails.ID)

		if err != nil {
			return nil, fmt.Errorf("get orders: %w", err)
		}

		orderDetails.Addons = addons

		orderDetailsList = append(orderDetailsList, orderDetails)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("get orders: %w", rows.Err())
	}

	return orderDetailsList, nil
}

func (orderSerice *OrderService) GetById(id string) (OrderDetails, error) {
	row := orderSerice.db.QueryRow(`SELECT
		orders.id,
		orders.customer_name,
		orders.customer_email,
		orders.phone_number,
		orders.no_of_persons,
		orders.total_price,
		orders.order_date,
		orders.ordered_at,
		theatres.id,
		theatres."name" ,
		theatres.description ,
		theatres.price ,
		theatres.additional_price_per_head ,
		theatres.max_capacity ,
		theatres.min_capacity ,
		theatres.default_capacity,
		theatres.created_at,
		theatres.updated_at,
		theatres.created_by,
		theatres.updated_by
		slots.id ,
		slots.start_time ,
		slots.end_time,
		slots.created_at,
		slots.updated_at,
		slots.created_by,
		slots.updated_by
		payments.razorpay_order_id,
		payments.razorpay_payment_id,
		payments.razorpay_signature,
		payments.status
	FROM
		orders
	JOIN theatres ON
		orders.theatre_id = theatres.id
	JOIN slots ON
		slots.id = orders.slot_id
	WHERE orders.id=$1;`, id)

	var orderDetails OrderDetails
	err := row.Scan(&orderDetails.ID, &orderDetails.CustomerName, &orderDetails.CustomerEmail, &orderDetails.PhoneNumber, &orderDetails.NoOfPersons, &orderDetails.TotalPrice, &orderDetails.OrderDate, &orderDetails.OrderedAt, &orderDetails.Theatre.ID, &orderDetails.Theatre.Name, &orderDetails.Theatre.Description, &orderDetails.Theatre.Price, &orderDetails.Theatre.AdditionalPricePerHead, &orderDetails.Theatre.MaxCapacity, &orderDetails.Theatre.MinCapacity, &orderDetails.Theatre.DefaultCapacity, &orderDetails.Theatre.CreatedAt, &orderDetails.Theatre.UpdatedAt, &orderDetails.Theatre.CreatedBy, &orderDetails.Theatre.UpdatedBy, &orderDetails.Slot.ID, &orderDetails.Slot.StartTime, &orderDetails.Slot.EndTime, &orderDetails.Slot.CreatedAt, &orderDetails.Slot.UpdatedAt, &orderDetails.Slot.CreatedBy, &orderDetails.Slot.UpdatedBy, &orderDetails.PaymentDetails.RazorpayOrderId, &orderDetails.PaymentDetails.RazorpayPaymentId, &orderDetails.PaymentDetails.RazorpaySignature, &orderDetails.PaymentDetails.Status)

	if err != nil {
		return OrderDetails{}, fmt.Errorf("get orders: %w", err)
	}

	addons, err := orderSerice.getAddonsForOrder(orderDetails.ID)

	if err != nil {
		return OrderDetails{}, fmt.Errorf("get orders: %w", err)
	}

	orderDetails.Addons = addons

	return orderDetails, nil
}

func (orderService OrderService) getAddonsForOrder(id string) ([]OrderAddonDetails, error) {
	rows, err := orderService.db.Query(`SELECT
		addons.id,
		addons.name,
		addons.category,
		addons.meta_data,
		addons.price,
		addons.created_at,
		addons.updated_at,
		addons.created_by,
		addons.updated_by,
		order_addons.quantity
	FROM
		order_addons
	JOIN addons ON
		order_addons.addon_id = addons.id
	WHERE order_addons.order_id = $1;`, id)

	if err != nil {
		return nil, fmt.Errorf("get orders: %w", err)
	}
	defer rows.Close()

	addons := make([]OrderAddonDetails, 0, 5)
	for rows.Next() {
		var addonDetails OrderAddonDetails
		err := rows.Scan(&addonDetails.ID, &addonDetails.Name, &addonDetails.Category, &addonDetails.MetaData, &addonDetails.Price, &addonDetails.CreatedAt, &addonDetails.UpdatedAt, &addonDetails.CreatedBy, &addonDetails.UpdatedBy, &addonDetails.Quantity)

		if err != nil {
			return nil, fmt.Errorf("get orders: %w", err)
		}
		addons = append(addons, addonDetails)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("get orders: %w", rows.Err())
	}

	return addons, err
}
