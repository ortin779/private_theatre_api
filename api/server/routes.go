package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ortin779/private_theatre_api/api/handlers"
	"github.com/ortin779/private_theatre_api/api/middleware"
	"github.com/ortin779/private_theatre_api/api/repository"
	"github.com/ortin779/private_theatre_api/api/service"
	"github.com/ortin779/private_theatre_api/config"
	"go.uber.org/zap"
)

func addRoutes(
	c *chi.Mux,
	logger *zap.Logger,
	db *sql.DB,
	cfg *config.Config,
) {

	// Repository initialization
	slotsRepository := repository.NewSlotsRepo(db)
	theatreRepository := repository.NewTheatreRepository(db)
	addonRepo := repository.NewAddonRepository(db)
	ordersRepo := repository.NewOrderRepository(db)
	usersRepo := repository.NewUsersRepository(db)
	paymentsRepo := repository.NewPaymentsRepository(db)

	// Service Initialization
	addonsService := service.NewAddonService(addonRepo)
	ordersService := service.NewOrdersService(ordersRepo)
	slotsService := service.NewSlotsService(slotsRepository)
	theatreService := service.NewTheatreService(theatreRepository)
	paymentService := service.NewRazorpayService(paymentsRepo, cfg.Razorpay)
	usersService := service.NewUsersService(usersRepo)

	// Handlers Initialization
	addonsHandler := handlers.NewAddonsHandler(logger, addonsService)
	authHandler := handlers.NewAuthHandler(logger, usersService)
	slotsHandler := handlers.NewSlotsHandler(logger, slotsService)
	ordersHandler := handlers.NewOrdersHandler(logger, ordersService, paymentService)
	paymentsHandler := handlers.NewPaymentHandler(logger, paymentService)
	theatreHandler := handlers.NewTheatreHandler(logger, theatreService)
	usersHandler := handlers.NewUsersHandler(logger, usersService)

	//add middlewares
	c.Use(middleware.RequestIdMiddleware)
	loggerMiddleware := middleware.LoggerMiddleware(logger)
	c.Use(loggerMiddleware)

	c.Get("/healthz", healthHandler)

	c.Post("/slots", middleware.AdminAuthorization(slotsHandler.HandleCreateSlot()))
	c.Get("/slots", slotsHandler.HandleSlotsGet())

	c.Post("/theatres", middleware.AdminAuthorization(theatreHandler.HandleCreateTheatre()))
	c.Get("/theatres", theatreHandler.HandleGetTheatres())
	c.Get("/theatres/{id}", theatreHandler.HandleGetTheatreDetails())

	c.Post("/addons", middleware.AdminAuthorization(addonsHandler.HandleCreateAddon()))
	c.Get("/addons", addonsHandler.HandleGetAddons())
	c.Get("/addons/categories", addonsHandler.HandleGetAddonCategories())

	c.Post("/orders", ordersHandler.HandleCreateOrder())
	c.Get("/orders", ordersHandler.HandleGetAllOrders())
	c.Get("/orders/{orderId}", ordersHandler.HandleGetOrderById())

	c.Post("/users", middleware.AdminAuthorization(usersHandler.HandleCreateUser()))

	c.Post("/login", authHandler.Login())
	c.Post("/refresh-token", authHandler.RefreshToken())

	c.Post("/verify-payment", paymentsHandler.VerifyPayment())

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
