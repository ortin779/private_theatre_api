package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ortin779/private_theatre_api/auth"
	"github.com/ortin779/private_theatre_api/models"
)

func Login(userStore models.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginParams models.LoginParams

		err := json.NewDecoder(r.Body).Decode(&loginParams)

		if err != nil {
			log.Println("invalid login params", loginParams)
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errs := loginParams.Validate(); len(errs) > 0 {
			log.Println("invalid login params", loginParams)
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		user, err := userStore.GetByEmail(loginParams.Email)

		if err != nil {
			log.Println(err.Error())
			if errors.Is(err, models.ErrNoUserWithEmail) {
				RespondWithError(w, http.StatusNotFound, fmt.Sprintf("no user found with given email: %s", loginParams.Email))
				return
			}
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		isValidPassword := auth.ComparePasswordToHash(user.Password, loginParams.Password)
		if !isValidPassword {
			log.Println("Authentication error, invalid credentials", loginParams)
			RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		accessToken, err := auth.GenerateAccessToken(user.ID, user.Roles)
		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Roles)
		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, models.LoginResponse{
			Token:        accessToken,
			RefreshToken: refreshToken,
		})

	}
}

func RefreshToken(userStore models.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshBody struct {
			RefreshToken string `json:"refresh_token"`
		}

		err := json.NewDecoder(r.Body).Decode(&refreshBody)

		if err != nil {
			log.Println("invalid login params", refreshBody)
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		claims, err := auth.ValidateToken(refreshBody.RefreshToken)

		if err != nil {
			log.Println(err)
			if errors.Is(err, auth.ErrTokenExpiry) {
				RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		token, err := auth.GenerateAccessToken(claims.UserId, claims.Roles)
		if err != nil {
			log.Println(err)
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusOK, models.LoginResponse{
			Token:        token,
			RefreshToken: refreshBody.RefreshToken,
		})
	}
}
