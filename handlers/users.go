package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/auth"
	"github.com/ortin779/private_theatre_api/models"
)

func HandleCreateUser(userStore models.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userParams models.UserParams

		err := json.NewDecoder(r.Body).Decode(&userParams)

		if err != nil {
			log.Println("error while decoding the create user request params", userParams)
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errs := userParams.Validate(); len(errs) > 0 {
			log.Println("user params validataion error", userParams)
			RespondWithJson(w, http.StatusBadRequest, errs)
			return
		}

		hashedPassword, err := auth.HashPassword(userParams.Password)

		if err != nil {
			log.Println("user params validataion error", err)
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

		err = userStore.Create(user)

		if err != nil {
			log.Println(err.Error())
			RespondWithError(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		RespondWithJson(w, http.StatusCreated, user)
	}
}
