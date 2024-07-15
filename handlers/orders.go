package handlers

import (
	"net/http"

	"github.com/ortin779/private_theatre_api/models"
)

func HandleCreateOrder(orderStore models.OrderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func HandleGetAllOrders(orderStore models.OrderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func HandleGetOrderById(orderStore models.OrderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
