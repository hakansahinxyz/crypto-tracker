// main.go
package main

import (
	"log"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/db"
	repository "github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/routes"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/services"
)

func main() {
	db.ConnectDatabase()

	walletBalanceRepo := repository.NewWalletBalanceRepository(db.DB)

	walletBalanceService := services.NewWalletBalanceService(walletBalanceRepo)

	router := routes.SetupRouter(walletBalanceService)

	cronService := services.NewCronService(walletBalanceService)
	cronService.StartCronJobs()

	port := "8080"
	log.Printf("Starting server on port %s...", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
