package middleware

import (
	"net/http"

	"github.com/ortin779/private_theatre_api/api/ctx"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	loggerMid := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqLogger := logger.With(
				zap.String("req-id", ctx.GetRequestId(r.Context())),
				zap.String("path", r.URL.Path),
				zap.String("method", r.Method),
				zap.String("addrs", r.RemoteAddr),
			)

			reqLogger.Info("incoming request")

			next.ServeHTTP(w, r)

			reqLogger.Info("request completed")
		}
		return http.HandlerFunc(fn)
	}

	return loggerMid
}
