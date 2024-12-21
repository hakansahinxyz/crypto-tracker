// main.go
package main

import (
	"log"

	"github.com/hakansahinxyz/crypto-tracker-backend/db"
	"github.com/hakansahinxyz/crypto-tracker-backend/routes"
)

func main() {
	// MySQL'e bağlantı kuruyoruz
	db.ConnectDatabase()

	router := routes.SetupRouter()

	port := "8080"
	log.Printf("Starting server on port %s...", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
