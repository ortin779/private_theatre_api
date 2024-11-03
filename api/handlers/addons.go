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

func HandleCreateAddon(addonsService service.AddonsService) http.HandlerFunc {
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

		userId, err := ctx.UserIdValue(r.Context())
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		addon := models.Addon{
			ID:        uuid.NewString(),
			Name:      addonParams.Name,
			Category:  addonParams.Category,
			Price:     addonParams.Price,
			MetaData:  addonParams.MetaData,
			CreatedBy: userId,
			UpdatedBy: userId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = addonsService.CreateAddon(addon)
		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, addon)
	}
}

func HandleGetAddons(addonService service.AddonsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addons, err := addonService.GetAllAddons()
		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, addons)
	}
}

func HandleGetAddonCategories(addonService service.AddonsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories := addonService.GetCategories()
		RespondWithJson(w, http.StatusCreated, categories)
	}
}
