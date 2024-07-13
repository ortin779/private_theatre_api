package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/models"
)

func addRoutes(c *chi.Mux, slotsStore models.SlotStore) {
	c.Get("/healthz", healthHandler)

	c.Get("/slots", handleSlotsGet(slotsStore))
	c.Post("/slots", handleCreateSlot(slotsStore))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleSlotsGet(slotsStore models.SlotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slots, err := slotsStore.GetSlots()

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, slots)
	}
}

func handleCreateSlot(slotsStore models.SlotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createSlotParams models.CreateSlotParams

		err := json.NewDecoder(r.Body).Decode(&createSlotParams)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		errs := createSlotParams.Validate()
		if len(errs) > 0 {
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		slot := models.Slot{
			ID:        uuid.New().String(),
			StartTime: convertMinutesToTime(createSlotParams.StartTime),
			EndTime:   convertMinutesToTime(createSlotParams.EndTime),
		}

		err = slotsStore.AddSlot(slot)
		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, slot)
	}
}

// a function to convert the given minutes to timestamp. It adds the given number of minutes from midnight
func convertMinutesToTime(minutes int) time.Time {
	// Get the current date at midnight
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Add the minutes to midnight
	return midnight.Add(time.Duration(minutes) * time.Minute)
}
