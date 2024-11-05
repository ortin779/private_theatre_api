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

func HandleCreateTheatre(ts service.TheatresService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.GetLogger(r.Context())
		var createTheatreParams models.CreateTheatreParams

		err := json.NewDecoder(r.Body).Decode(&createTheatreParams)

		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if errs := createTheatreParams.Validate(); len(errs) > 0 {
			logger.Error("invalid request", zap.Any("errors", errs))
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		userId, err := ctx.UserIdValue(r.Context())
		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
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

		err = ts.Create(theatre, createTheatreParams.Slots)

		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, theatre)

	}
}

func HandleGetTheatres(ts service.TheatresService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.GetLogger(r.Context())
		theatres, err := ts.GetTheatres()

		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusOK, theatres)
	}
}

func HandleGetTheatreDetails(ts service.TheatresService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.GetLogger(r.Context())
		id := r.PathValue("id")
		if id == "" {
			logger.Error("invalid request", zap.String("error", "theatre id can not be empty"))
			RespondWithError(w, http.StatusBadRequest, "invalid theatre id")
		}

		theatres, err := ts.GetTheatreDetails(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("not found", zap.String("error", "no theatre found with given id"))
				RespondWithError(w, 404, "no theatre found with given details")
				return
			}
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusOK, theatres)
	}
}
