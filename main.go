package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/OPGLOL/opgl-data/internal/api"
	"github.com/OPGLOL/opgl-data/internal/config"
	"github.com/OPGLOL/opgl-data/internal/services"
)

func main() {
	// Load configuration
	configuration := config.LoadConfig()

	// Initialize Riot service
	riotService := services.NewRiotService(configuration.RiotAPIKey)

	// Initialize HTTP handler
	handler := api.NewHandler(riotService)

	// Set up router
	router := api.SetupRouter(handler)

	// Start server
	serverAddress := fmt.Sprintf(":%s", configuration.ServerPort)
	log.Printf("OPGL Data Service starting on port %s", configuration.ServerPort)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
