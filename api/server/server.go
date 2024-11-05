package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ortin779/private_theatre_api/config"
	"go.uber.org/zap"
)

func NewServer(
	logger *zap.Logger,
	db *sql.DB,
	cfg *config.Config,
) http.Handler {
	router := chi.NewRouter()

	addRoutes(router, logger, db, cfg)

	return router
}
