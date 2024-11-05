package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ortin779/private_theatre_api/api/ctx"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				reqLogger := logger.With(
					zap.String("req-id", uuid.NewString()),
					zap.String("path", r.URL.Path),
				)

				ctx := ctx.WithLogger(r.Context(), reqLogger)
				r = r.WithContext(ctx)

				reqLogger.Info("incoming request")

				next.ServeHTTP(w, r)

				reqLogger.Info("request completed")
			},
		)
	}
}
