package middleware

import (
	"net/http"

	"github.com/ortin779/private_theatre_api/api/ctx"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := ctx.WithLogger(r.Context(), logger)
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
			},
		)
	}
}
