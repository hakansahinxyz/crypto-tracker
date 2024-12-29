// main.go
package main

import (
	"log"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/db"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/routes"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/services"
)

func main() {
	db.ConnectDatabase()

	//go services.StartPriceService()
	go services.StartBalanceService()

	router := routes.SetupRouter()

	port := "8080"
	log.Printf("Starting server on port %s...", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
