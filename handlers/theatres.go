package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
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

		theatre := models.Theatre{
			ID:                     uuid.NewString(),
			Name:                   createTheatreParams.Name,
			Description:            createTheatreParams.Description,
			Price:                  createTheatreParams.Price,
			AdditionalPricePerHead: createTheatreParams.AdditionalPricePerHead,
			MaxCapacity:            createTheatreParams.MaxCapacity,
			MinCapacity:            createTheatreParams.MinCapacity,
			DefaultCapacity:        createTheatreParams.DefaultCapacity,
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
