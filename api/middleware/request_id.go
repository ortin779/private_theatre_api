package middleware

import (
	"net/http"

	"github.com/ortin779/private_theatre_api/api/ctx"
)

func RequestIdMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctxWithReqId := ctx.WithRequestId(r.Context())

		r = r.WithContext(ctxWithReqId)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
