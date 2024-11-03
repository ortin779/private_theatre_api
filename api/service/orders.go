package service

import (
	"database/sql"
	"errors"

	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/repository"
)

type OrdersService struct {
	ordersRepo repository.OrdersRepository
}

func NewOrdersService(ordersRepo repository.OrdersRepository) OrdersService {
	return OrdersService{
		ordersRepo: ordersRepo,
	}
}

var ErrDuplicateOrder = errors.New("order already exists for the given theatre and slot")

func (o *OrdersService) Create(order models.Order) error {
	_, err := o.ordersRepo.GetOrderByTheatreIdAndSlotIdAndOrderDate(order.SlotId, order.TheatreId, order.OrderDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDuplicateOrder
		}
		return err
	}

	return o.ordersRepo.Create(order)
}

func (o *OrdersService) GetAll() ([]models.OrderDetails, error) {
	return o.ordersRepo.GetAll()
}

func (o *OrdersService) GetById(id string) (*models.OrderDetails, error) {
	return o.ordersRepo.GetById(id)
}
