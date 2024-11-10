package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/api/auth"
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/service"
	"go.uber.org/zap"
)

type UsersHandler struct {
	usersService service.UsersService
	logger       *zap.Logger
}

func NewUsersHandler(logger *zap.Logger, usersService service.UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
		logger:       logger,
	}
}

func (usrHandler *UsersHandler) HandleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userParams models.UserParams

		err := json.NewDecoder(r.Body).Decode(&userParams)

		if err != nil {
			usrHandler.logger.Error("invalid request", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errs := userParams.Validate(); len(errs) > 0 {
			usrHandler.logger.Error("invalid request", zap.Any("errors", errs))
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		hashedPassword, err := auth.HashPassword(userParams.Password)

		if err != nil {
			usrHandler.logger.Error("internal server error", zap.String("error", err.Error()))
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

		err = usrHandler.usersService.Create(user)

		if err != nil {
			usrHandler.logger.Error("internal server error", zap.String("error", err.Error()))
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusCreated, user)
	}
}
