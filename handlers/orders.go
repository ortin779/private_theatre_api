package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/models"
)

func HandleCreateOrder(orderStore models.OrderStore, paymentService *models.RazorpayService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var orderParams models.OrderParams

		err := json.NewDecoder(r.Body).Decode(&orderParams)

		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		if errs := orderParams.Validate(); len(errs) > 0 {
			log.Println(err)
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		normalizedPrice := orderParams.TotalPrice * 100

		razorpayOrderId, err := paymentService.CreateOrder(normalizedPrice)

		if err != nil {
			log.Println(err)
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

		err = orderStore.Create(order)

		if err != nil {
			log.Println(err)
			if errors.Is(err, models.ErrDuplicateOrder) {
				RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}
		RespondWithJson(w, http.StatusCreated, order)
	}
}

func HandleGetAllOrders(orderStore models.OrderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := orderStore.GetAll()

		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, orders)
	}
}

func HandleGetOrderById(orderStore models.OrderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderId := r.PathValue("orderId")
		if _, err := uuid.Parse(orderId); err != nil {
			RespondWithError(w, http.StatusNotFound, "invalid order id")
			return
		}

		orderDetails, err := orderStore.GetById(orderId)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, orderDetails)
	}
}
