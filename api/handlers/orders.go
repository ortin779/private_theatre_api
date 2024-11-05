package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/api/ctx"
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/service"
	"go.uber.org/zap"
)

func HandleCreateOrder(ordersService service.OrdersService, paymentService service.RazorpayService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.GetLogger(r.Context())

		var orderParams models.OrderParams

		err := json.NewDecoder(r.Body).Decode(&orderParams)

		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		if errs := orderParams.Validate(); len(errs) > 0 {
			logger.Error("bad request", zap.Any("errs", errs))
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		normalizedPrice := orderParams.TotalPrice * 100

		razorpayOrderId, err := paymentService.CreateOrder(normalizedPrice)

		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
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

		err = ordersService.Create(order)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("bad request", zap.String("error", err.Error()))
				RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
		RespondWithJson(w, http.StatusCreated, order)
	}
}

func HandleGetAllOrders(ordersService service.OrdersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.GetLogger(r.Context())
		orders, err := ordersService.GetAll()

		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, orders)
	}
}

func HandleGetOrderById(ordersService service.OrdersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.GetLogger(r.Context())
		orderId := r.PathValue("orderId")
		if _, err := uuid.Parse(orderId); err != nil {
			logger.Error("not found", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusNotFound, "invalid order id")
			return
		}

		orderDetails, err := ordersService.GetById(orderId)
		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, orderDetails)
	}
}
