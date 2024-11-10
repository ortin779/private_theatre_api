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

type TheatreHandler struct {
	logger         *zap.Logger
	theatreService service.TheatresService
}

func NewTheatreHandler(logger *zap.Logger, theatreService service.TheatresService) *TheatreHandler {
	return &TheatreHandler{
		logger:         logger,
		theatreService: theatreService,
	}
}

func (thrHandler *TheatreHandler) HandleCreateTheatre() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var createTheatreParams models.CreateTheatreParams

		err := json.NewDecoder(r.Body).Decode(&createTheatreParams)

		if err != nil {
			thrHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if errs := createTheatreParams.Validate(); len(errs) > 0 {
			thrHandler.logger.Error("invalid request", zap.Any("errors", errs))
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		userId, err := ctx.UserIdValue(r.Context())
		if err != nil {
			thrHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		theatre := models.Theatre{
			ID:                     uuid.NewString(),
			Name:                   createTheatreParams.Name,
			Description:            createTheatreParams.Description,
			Price:                  createTheatreParams.Price,
			AdditionalPricePerHead: createTheatreParams.AdditionalPricePerHead,
			MaxCapacity:            createTheatreParams.MaxCapacity,
			MinCapacity:            createTheatreParams.MinCapacity,
			DefaultCapacity:        createTheatreParams.DefaultCapacity,
			CreatedBy:              userId,
			UpdatedBy:              userId,
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
		}

		err = thrHandler.theatreService.Create(theatre, createTheatreParams.Slots)

		if err != nil {
			thrHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, theatre)

	}
}

func (thrHandler *TheatreHandler) HandleGetTheatres() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		theatres, err := thrHandler.theatreService.GetTheatres()

		if err != nil {
			thrHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusOK, theatres)
	}
}

func (thrHandler *TheatreHandler) HandleGetTheatreDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id")
		if id == "" {
			thrHandler.logger.Error("invalid request", zap.String("error", "theatre id can not be empty"))
			RespondWithError(w, http.StatusBadRequest, "invalid theatre id")
		}

		theatres, err := thrHandler.theatreService.GetTheatreDetails(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				thrHandler.logger.Error("not found", zap.String("error", "no theatre found with given id"))
				RespondWithError(w, 404, "no theatre found with given details")
				return
			}
			thrHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusOK, theatres)
	}
}
