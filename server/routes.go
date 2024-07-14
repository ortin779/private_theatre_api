package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ortin779/private_theatre_api/handlers"
	"github.com/ortin779/private_theatre_api/models"
)

func addRoutes(c *chi.Mux, slotsStore models.SlotStore) {
	c.Get("/healthz", healthHandler)

	c.Get("/slots", handlers.HandleSlotsGet(slotsStore))
	c.Post("/slots", handlers.HandleCreateSlot(slotsStore))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
