package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ortin779/private_theatre_api/api/repository"
	"github.com/ortin779/private_theatre_api/api/server"
	"github.com/ortin779/private_theatre_api/api/service"
	"github.com/ortin779/private_theatre_api/config"
)

func main() {
	cfg, err := config.LoadConfigFromEnv()

	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := cfg.Postgres.Open()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

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

	svr := server.NewServer(slotsService, theatreService, addonsService, ordersService, usersService, paymentService)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler: svr,
	}

	fmt.Println("Server stared on port ", cfg.Server.Port)
	err = httpServer.ListenAndServe()

	if err != nil {
		log.Fatalln(err.Error())
	}
}
