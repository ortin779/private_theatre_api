package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/api/auth"
	"github.com/ortin779/private_theatre_api/api/ctx"
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/service"
	"go.uber.org/zap"
)

func HandleCreateUser(usersService service.UsersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := ctx.GetLogger(r.Context())
		var userParams models.UserParams

		err := json.NewDecoder(r.Body).Decode(&userParams)

		if err != nil {
			logger.Error("invalid request", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errs := userParams.Validate(); len(errs) > 0 {
			logger.Error("invalid request", zap.Any("errors", errs))
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		hashedPassword, err := auth.HashPassword(userParams.Password)

		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithJson(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		user := models.User{
			ID:       uuid.NewString(),
			Name:     userParams.Name,
			Email:    strings.ToLower(userParams.Email),
			Password: hashedPassword,
			Roles:    userParams.Roles,
		}

		err = usersService.Create(user)

		if err != nil {
			logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusCreated, user)
	}
}
