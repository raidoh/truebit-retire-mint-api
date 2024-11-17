package main

import (
	"log"
	"net/http"
	"os"

	"truebit-api/internal/api"
	"truebit-api/internal/client"
	"truebit-api/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Ethereum client
	ethClient, err := client.NewEthereumClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize ethereum client: %v", err)
	}

	// Set up router
	router := api.NewRouter(ethClient)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on http://localhost:%s/truebit", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
