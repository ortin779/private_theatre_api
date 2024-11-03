package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ortin779/private_theatre_api/api/models"
)

type OrdersRepository interface {
	Create(order models.Order) error
	GetAll() ([]models.OrderDetails, error)
	GetById(id string) (*models.OrderDetails, error)
	GetOrderByTheatreIdAndSlotIdAndOrderDate(slotId, theatreId string, orderDate time.Time) (*models.OrderDetails, error)
}

type ordersRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (ordersRepo *ordersRepository) GetOrderByTheatreIdAndSlotIdAndOrderDate(slotId, theatreId string, orderDate time.Time) (*models.OrderDetails, error) {
	var orderId string
	row := ordersRepo.db.QueryRow(`SELECT id FROM orders
        WHERE theatre_id = $1 AND slot_id = $2 AND order_date = $3;
    `, theatreId, slotId, orderDate.Format(time.DateOnly))

	err := row.Scan(&orderId)
	if err != nil {
		return nil, err
	}

	return ordersRepo.GetById(orderId)
}

func (ordersRepo *ordersRepository) Create(order models.Order) error {
	tx, err := ordersRepo.db.Begin()

	if err != nil {
		return fmt.Errorf("create order: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(`INSERT INTO orders(
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

func (ordersRepo *ordersRepository) GetAll() ([]models.OrderDetails, error) {

	rows, err := ordersRepo.db.Query(`SELECT
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

	orderDetailsList := make([]models.OrderDetails, 0, 5)
	for rows.Next() {
		var orderDetails models.OrderDetails
		err := rows.Scan(&orderDetails.ID, &orderDetails.CustomerName, &orderDetails.CustomerEmail, &orderDetails.PhoneNumber, &orderDetails.NoOfPersons, &orderDetails.TotalPrice, &orderDetails.OrderDate, &orderDetails.OrderedAt, &orderDetails.Theatre.ID, &orderDetails.Theatre.Name, &orderDetails.Theatre.Description, &orderDetails.Theatre.Price, &orderDetails.Theatre.AdditionalPricePerHead, &orderDetails.Theatre.MaxCapacity, &orderDetails.Theatre.MinCapacity, &orderDetails.Theatre.DefaultCapacity, &orderDetails.Slot.ID, &orderDetails.Slot.StartTime, &orderDetails.Slot.EndTime, &orderDetails.PaymentDetails.RazorpayOrderId, &orderDetails.PaymentDetails.RazorpayPaymentId, &orderDetails.PaymentDetails.RazorpaySignature, &orderDetails.PaymentDetails.Status)

		if err != nil {
			return nil, fmt.Errorf("get orders: %w", err)
		}

		addons, err := ordersRepo.getAddonsForOrder(orderDetails.ID)

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

func (ordersRepo *ordersRepository) GetById(id string) (*models.OrderDetails, error) {
	row := ordersRepo.db.QueryRow(`SELECT
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

	var orderDetails models.OrderDetails
	err := row.Scan(&orderDetails.ID, &orderDetails.CustomerName, &orderDetails.CustomerEmail, &orderDetails.PhoneNumber, &orderDetails.NoOfPersons, &orderDetails.TotalPrice, &orderDetails.OrderDate, &orderDetails.OrderedAt, &orderDetails.Theatre.ID, &orderDetails.Theatre.Name, &orderDetails.Theatre.Description, &orderDetails.Theatre.Price, &orderDetails.Theatre.AdditionalPricePerHead, &orderDetails.Theatre.MaxCapacity, &orderDetails.Theatre.MinCapacity, &orderDetails.Theatre.DefaultCapacity, &orderDetails.Theatre.CreatedAt, &orderDetails.Theatre.UpdatedAt, &orderDetails.Theatre.CreatedBy, &orderDetails.Theatre.UpdatedBy, &orderDetails.Slot.ID, &orderDetails.Slot.StartTime, &orderDetails.Slot.EndTime, &orderDetails.Slot.CreatedAt, &orderDetails.Slot.UpdatedAt, &orderDetails.Slot.CreatedBy, &orderDetails.Slot.UpdatedBy, &orderDetails.PaymentDetails.RazorpayOrderId, &orderDetails.PaymentDetails.RazorpayPaymentId, &orderDetails.PaymentDetails.RazorpaySignature, &orderDetails.PaymentDetails.Status)

	if err != nil {
		return nil, fmt.Errorf("get orders: %w", err)
	}

	addons, err := ordersRepo.getAddonsForOrder(orderDetails.ID)

	if err != nil {
		return nil, fmt.Errorf("get orders: %w", err)
	}

	orderDetails.Addons = addons

	return &orderDetails, nil
}

func (ordersRepo *ordersRepository) getAddonsForOrder(id string) ([]models.OrderAddonDetails, error) {
	rows, err := ordersRepo.db.Query(`SELECT
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

	addons := make([]models.OrderAddonDetails, 0, 5)
	for rows.Next() {
		var addonDetails models.OrderAddonDetails
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
