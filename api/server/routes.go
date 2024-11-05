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

	//add middlewares
	loggerMiddleware := middleware.LoggerMiddleware(logger)
	c.Use(loggerMiddleware)

	c.Get("/healthz", healthHandler)

	c.Post("/slots", middleware.AdminAuthorization(handlers.HandleCreateSlot(slotsService)))
	c.Get("/slots", handlers.HandleSlotsGet(slotsService))

	c.Post("/theatres", middleware.AdminAuthorization(handlers.HandleCreateTheatre(theatreService)))
	c.Get("/theatres", handlers.HandleGetTheatres(theatreService))
	c.Get("/theatres/{id}", handlers.HandleGetTheatreDetails(theatreService))

	c.Post("/addons", middleware.AdminAuthorization(handlers.HandleCreateAddon(addonsService)))
	c.Get("/addons", handlers.HandleGetAddons(addonsService))
	c.Get("/addons/categories", handlers.HandleGetAddonCategories(addonsService))

	c.Post("/orders", handlers.HandleCreateOrder(ordersService, paymentService))
	c.Get("/orders", handlers.HandleGetAllOrders(ordersService))
	c.Get("/orders/{orderId}", handlers.HandleGetOrderById(ordersService))

	c.Post("/users", middleware.AdminAuthorization(handlers.HandleCreateUser(usersService)))

	c.Post("/login", handlers.Login(usersService))
	c.Post("/refresh-token", handlers.RefreshToken(usersService))

	c.Post("/verify-payment", handlers.VerifyPayment(paymentService))

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
