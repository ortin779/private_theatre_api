package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ortin779/private_theatre_api/models"
)

func NewServer(
	slotsStore models.SlotStore,
	theatreStore models.TheatreStore,
	addonStore models.AddonStore,
	orderStore models.OrderStore,
	userStore models.UserStore,
) http.Handler {
	router := chi.NewRouter()

	addRoutes(router, slotsStore, theatreStore, addonStore, orderStore, userStore)

	return router
}
