package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/api/ctx"
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/service"
	"go.uber.org/zap"
)

type SlotsHandler struct {
	logger       *zap.Logger
	slotsService service.SlotsService
}

func NewSlotsHandler(logger *zap.Logger, slotsService service.SlotsService) *SlotsHandler {
	return &SlotsHandler{
		logger:       logger,
		slotsService: slotsService,
	}
}

func (slotsHandler *SlotsHandler) HandleSlotsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slots, err := slotsHandler.slotsService.GetSlots()

		if err != nil {
			slotsHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, slots)
	}
}

func (slotsHandler *SlotsHandler) HandleCreateSlot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createSlotParams models.CreateSlotParams

		err := json.NewDecoder(r.Body).Decode(&createSlotParams)

		if err != nil {
			slotsHandler.logger.Error("internal server error", zap.String("error", err.Error()))
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
			slotsHandler.logger.Error("internal server error", zap.String("error", err.Error()))
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

		err = slotsHandler.slotsService.AddSlot(slot)
		if err != nil {
			slotsHandler.logger.Error("internal server error", zap.String("error", err.Error()))
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
