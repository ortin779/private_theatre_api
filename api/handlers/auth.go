package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ortin779/private_theatre_api/api/auth"
	"github.com/ortin779/private_theatre_api/api/models"
	"github.com/ortin779/private_theatre_api/api/service"
	"go.uber.org/zap"
)

type AuthHandler struct {
	usersService service.UsersService
	logger       *zap.Logger
}

func NewAuthHandler(logger *zap.Logger, usersService service.UsersService) *AuthHandler {
	return &AuthHandler{
		usersService: usersService,
		logger:       logger,
	}
}

func (authHandler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginParams models.LoginParams

		err := json.NewDecoder(r.Body).Decode(&loginParams)

		if err != nil {
			authHandler.logger.Error("invalid login params ", zap.Any("error", err))
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errs := loginParams.Validate(); len(errs) > 0 {
			authHandler.logger.Error("invalid login params", zap.Any("errors", errs))
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		user, err := authHandler.usersService.GetByEmail(loginParams.Email)

		if err != nil {
			authHandler.logger.Error(err.Error())
			if errors.Is(err, models.ErrNoUserWithEmail) {
				RespondWithError(w, http.StatusNotFound, fmt.Sprintf("no user found with given email: %s", loginParams.Email))
				return
			}
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		isValidPassword := auth.ComparePasswordToHash(user.Password, loginParams.Password)
		if !isValidPassword {
			authHandler.logger.Error("Authentication error, invalid credentials", zap.Any("error", loginParams))
			RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		accessToken, err := auth.GenerateAccessToken(user.ID, user.Roles)
		if err != nil {
			authHandler.logger.Error(err.Error())
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Roles)
		if err != nil {
			authHandler.logger.Error(err.Error())
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, models.LoginResponse{
			Token:        accessToken,
			RefreshToken: refreshToken,
		})

	}
}

func (authHandler *AuthHandler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshBody struct {
			RefreshToken string `json:"refresh_token"`
		}

		err := json.NewDecoder(r.Body).Decode(&refreshBody)

		if err != nil {
			authHandler.logger.Error("invalid login params")
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		claims, err := auth.ValidateToken(refreshBody.RefreshToken)

		if err != nil {
			authHandler.logger.Error(err.Error())
			if errors.Is(err, auth.ErrTokenExpiry) {
				RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		token, err := auth.GenerateAccessToken(claims.UserId, claims.Roles)
		if err != nil {
			authHandler.logger.Error(err.Error())
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, models.LoginResponse{
			Token:        token,
			RefreshToken: refreshBody.RefreshToken,
		})
	}
}
