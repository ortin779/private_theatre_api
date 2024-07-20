package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/ctx"
	"github.com/ortin779/private_theatre_api/models"
)

func HandleCreateTheatre(ts models.TheatreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var createTheatreParams models.CreateTheatreParams

		err := json.NewDecoder(r.Body).Decode(&createTheatreParams)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if errs := createTheatreParams.Validate(); len(errs) > 0 {
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		userId, err := ctx.UserIdValue(r.Context())
		if err != nil {
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
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, theatre)

	}
}

func HandleGetTheatres(ts models.TheatreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		theatres, err := ts.GetTheatres()

		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusOK, theatres)
	}
}

func HandleGetTheatreDetails(ts models.TheatreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			RespondWithError(w, http.StatusBadRequest, "invalid theatre id")
		}

		theatres, err := ts.GetTheatreDetails(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				RespondWithError(w, 404, "no theatre found with given details")
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusOK, theatres)
	}
}
