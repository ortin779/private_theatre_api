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

func HandleCreateAddon(addonsService service.AddonsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.GetLogger(r.Context())
		var addonParams models.AddonParams

		err := json.NewDecoder(r.Body).Decode(&addonParams)

		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if errs := addonParams.Validate(); len(errs) > 0 {
			logger.Error("bad request", zap.Any("errors", errs))
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		userId, err := ctx.UserIdValue(r.Context())
		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
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
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, addon)
	}
}

func HandleGetAddons(addonService service.AddonsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.GetLogger(r.Context())
		addons, err := addonService.GetAllAddons()
		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
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
