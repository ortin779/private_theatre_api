package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ortin779/private_theatre_api/api/service"
	"go.uber.org/zap"
)

func NewServer(
	logger *zap.Logger,
	slotsService service.SlotsService,
	theatreService service.TheatresService,
	addonService service.AddonsService,
	ordersService service.OrdersService,
	usersService service.UsersService,
	razorpayService service.RazorpayService,
) http.Handler {
	router := chi.NewRouter()

	addRoutes(router, logger, slotsService, theatreService, addonService, ordersService, usersService, razorpayService)

	return router
}
