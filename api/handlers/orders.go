package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/service"
	"go.uber.org/zap"
)

type OrdersHandler struct {
	logger          *zap.Logger
	ordersService   service.OrdersService
	paymentsService service.RazorpayService
}

func NewOrdersHandler(logger *zap.Logger,
	ordersService service.OrdersService,
	paymentsService service.RazorpayService) *OrdersHandler {
	return &OrdersHandler{
		logger:          logger,
		ordersService:   ordersService,
		paymentsService: paymentsService,
	}
}

func (orderHandler *OrdersHandler) HandleCreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var orderParams models.OrderParams

		err := json.NewDecoder(r.Body).Decode(&orderParams)

		if err != nil {
			orderHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		if errs := orderParams.Validate(); len(errs) > 0 {
			orderHandler.logger.Error("bad request", zap.Any("errs", errs))
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		normalizedPrice := orderParams.TotalPrice * 100

		razorpayOrderId, err := orderHandler.paymentsService.CreateOrder(normalizedPrice)

		if err != nil {
			orderHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong, while creating payment")
			return
		}

		order := models.Order{
			ID:              uuid.NewString(),
			CustomerName:    orderParams.CustomerName,
			CustomerEmail:   orderParams.CustomerEmail,
			PhoneNumber:     orderParams.PhoneNumber,
			TheatreId:       orderParams.TheatreId,
			Addons:          orderParams.Addons,
			SlotId:          orderParams.SlotId,
			NoOfPersons:     orderParams.NoOfPersons,
			TotalPrice:      orderParams.TotalPrice,
			OrderDate:       orderParams.OrderDate,
			OrderedAt:       time.Now(),
			RazorpayOrderId: razorpayOrderId,
		}

		err = orderHandler.ordersService.Create(order)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				orderHandler.logger.Error("bad request", zap.String("error", err.Error()))
				RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			orderHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
		RespondWithJson(w, http.StatusCreated, order)
	}
}

func (orderHandler *OrdersHandler) HandleGetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		orders, err := orderHandler.ordersService.GetAll()

		if err != nil {
			orderHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, orders)
	}
}

func (orderHandler *OrdersHandler) HandleGetOrderById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		orderId := r.PathValue("orderId")
		if _, err := uuid.Parse(orderId); err != nil {
			orderHandler.logger.Error("not found", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusNotFound, "invalid order id")
			return
		}

		orderDetails, err := orderHandler.ordersService.GetById(orderId)
		if err != nil {
			orderHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, orderDetails)
	}
}
