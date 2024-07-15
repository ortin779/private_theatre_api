package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ortin779/private_theatre_api/handlers"
	"github.com/ortin779/private_theatre_api/models"
)

func addRoutes(
	c *chi.Mux,
	slotsStore models.SlotStore,
	theatreStore models.TheatreStore,
	addonStore models.AddonStore,
	orderStore models.OrderStore,
) {
	c.Get("/healthz", healthHandler)

	c.Get("/slots", handlers.HandleSlotsGet(slotsStore))
	c.Post("/slots", handlers.HandleCreateSlot(slotsStore))

	c.Post("/theatres", handlers.HandleCreateTheatre(theatreStore))
	c.Get("/theatres", handlers.HandleGetTheatres(theatreStore))
	c.Get("/theatres/{id}", handlers.HandleGetTheatreDetails(theatreStore))

	c.Post("/addons", handlers.HandleCreateAddon(addonStore))
	c.Get("/addons", handlers.HandleGetAddons(addonStore))
	c.Get("/addons/categories", handlers.HandleGetAddonCategories(addonStore))

	c.Post("/orders", handlers.HandleCreateOrder(orderStore))
	c.Get("/orders", handlers.HandleGetAllOrders(orderStore))
	c.Get("/orders/{orderId}", handlers.HandleGetOrderById(orderStore))

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
