package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ortin779/private_theatre_api/api/server"
	"github.com/ortin779/private_theatre_api/config"
	"github.com/ortin779/private_theatre_api/logger"
	"go.uber.org/zap"
)

func run(ctx context.Context, logger *zap.Logger) error {
	cfg, err := config.LoadConfigFromEnv()

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	db, err := cfg.Postgres.Open()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	defer db.Close()

	svr := server.NewServer(logger, db, cfg)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler: svr,
	}

	fmt.Println("Server stared on port ", cfg.Server.Port)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Info("shutdown", zap.String("status", "shutdown started"), zap.Any("signal", sig))
		defer logger.Info("shutdown", zap.String("status", "shutdown completed"), zap.Any("signal", sig))

		ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Web.ShutdownTimeout))
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func main() {

	// logger
	logger := logger.NewLogger()

	defer logger.Sync()

	ctx := context.Background()

	if err := run(ctx, logger); err != nil {
		logger.Error("startup", zap.String("msg", err.Error()))
		os.Exit(1)
	}
}
