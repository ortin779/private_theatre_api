package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/ortin779/private_theatre_api/auth"
	"github.com/ortin779/private_theatre_api/ctx"
	"github.com/ortin779/private_theatre_api/handlers"
)

func AdminAuthorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := getTokenFromReqest(r)
		claims, err := auth.ValidateToken(accessToken)

		if err != nil {
			handlers.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !slices.Contains(claims.Roles, "admin") {
			handlers.RespondWithError(w, http.StatusForbidden, "need admin previlizes to access")
			return
		}

		ctx := ctx.WithUserId(r.Context(), claims.UserId)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func getTokenFromReqest(r *http.Request) string {
	token := r.Header.Get("Authorization")

	tokenParts := strings.Split(token, " ")
	if len(tokenParts) == 2 {
		return tokenParts[1]
	}
	return ""
}
