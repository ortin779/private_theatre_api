package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/api/ctx"
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/service"
)

func HandleSlotsGet(slotsService service.SlotsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slots, err := slotsService.GetSlots()

		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, slots)
	}
}

func HandleCreateSlot(slotsService service.SlotsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createSlotParams models.CreateSlotParams

		err := json.NewDecoder(r.Body).Decode(&createSlotParams)

		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		errs := createSlotParams.Validate()
		if len(errs) > 0 {
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		userId, err := ctx.UserIdValue(r.Context())
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		slot := models.Slot{
			ID:        uuid.New().String(),
			StartTime: convertMinutesToTime(createSlotParams.StartTime),
			EndTime:   convertMinutesToTime(createSlotParams.EndTime),
			CreatedBy: userId,
			UpdatedBy: userId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = slotsService.AddSlot(slot)
		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusCreated, slot)
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
