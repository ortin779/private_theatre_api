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

type AddonsHandler struct {
	addonsService service.AddonsService
	logger        *zap.Logger
}

func NewAddonsHandler(logger *zap.Logger, addonsService service.AddonsService) *AddonsHandler {
	return &AddonsHandler{
		addonsService: addonsService,
		logger:        logger,
	}
}

func (ah *AddonsHandler) HandleCreateAddon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var addonParams models.AddonParams

		err := json.NewDecoder(r.Body).Decode(&addonParams)

		if err != nil {
			ah.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		if errs := addonParams.Validate(); len(errs) > 0 {
			ah.logger.Error("bad request", zap.Any("errors", errs))
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		userId, err := ctx.UserIdValue(r.Context())
		if err != nil {
			ah.logger.Error("internal server error", zap.String("error", err.Error()))
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

		err = ah.addonsService.CreateAddon(addon)
		if err != nil {
			ah.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, addon)
	}
}

func (ah *AddonsHandler) HandleGetAddons() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addons, err := ah.addonsService.GetAllAddons()
		if err != nil {
			ah.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		RespondWithJson(w, http.StatusCreated, addons)
	}
}

func (ah *AddonsHandler) HandleGetAddonCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories := ah.addonsService.GetCategories()
		RespondWithJson(w, http.StatusCreated, categories)
	}
}
