package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ortin779/private_theatre_api/api/handlers"
	"github.com/ortin779/private_theatre_api/api/middleware"
	"github.com/ortin779/private_theatre_api/api/service"
	"go.uber.org/zap"
)

func addRoutes(
	c *chi.Mux,
	logger *zap.Logger,
	slotsService service.SlotsService,
	theatresService service.TheatresService,
	addonsService service.AddonsService,
	ordersService service.OrdersService,
	usersService service.UsersService,
	paymentService service.RazorpayService,
) {
	loggerMiddleware := middleware.LoggerMiddleware(logger)
	c.Use(loggerMiddleware)

	c.Get("/healthz", healthHandler)

	c.Post("/slots", middleware.AdminAuthorization(handlers.HandleCreateSlot(slotsService)))
	c.Get("/slots", handlers.HandleSlotsGet(slotsService))

	c.Post("/theatres", middleware.AdminAuthorization(handlers.HandleCreateTheatre(theatresService)))
	c.Get("/theatres", handlers.HandleGetTheatres(theatresService))
	c.Get("/theatres/{id}", handlers.HandleGetTheatreDetails(theatresService))

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
