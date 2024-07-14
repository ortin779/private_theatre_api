package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ortin779/private_theatre_api/config"
	"github.com/ortin779/private_theatre_api/models"
	"github.com/ortin779/private_theatre_api/server"
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

	slotsStore := models.NewSlotService(db)
	theatreStore := models.NewTheatreService(db)
	addonStore := models.NewAddonStore(db)

	svr := server.NewServer(slotsStore, theatreStore, addonStore)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler: svr,
	}

	fmt.Println("Server stared on port ", cfg.Server.Port)
	err = httpServer.ListenAndServe()

	if err != nil {
		log.Fatalf(err.Error())
	}
}
