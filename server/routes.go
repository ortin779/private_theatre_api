package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ortin779/private_theatre_api/models"
)

func addRoutes(c *chi.Mux, slotsStore models.SlotStore) {
	c.Get("/healthz", healthHandler)

	c.Get("/slots", handleSlotsGet(slotsStore))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleSlotsGet(slotsStore models.SlotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slots, err := slotsStore.GetSlots()

		if err != nil {
			log.Println(err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(slots)

		if err != nil {
			log.Println(err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}
