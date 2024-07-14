package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/models"
)

func HandleCreateAddon(addonStore models.AddonStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var addonParams models.AddonParams

		err := json.NewDecoder(r.Body).Decode(&addonParams)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if errs := addonParams.Validate(); len(errs) > 0 {
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		addon := models.Addon{
			ID:       uuid.NewString(),
			Name:     addonParams.Name,
			Category: addonParams.Category,
			Price:    addonParams.Price,
			MetaData: addonParams.MetaData,
		}

		err = addonStore.Create(addon)
		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, addon)
	}
}

func HandleGetAddons(addonStore models.AddonStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addons, err := addonStore.GetAddons()
		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, addons)
	}
}

func HandleGetAddonCategories(addonStore models.AddonStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories := addonStore.GetCategories()
		RespondWithJson(w, http.StatusCreated, categories)
	}
}
