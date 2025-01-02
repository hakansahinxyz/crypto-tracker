// main.go
package main

import (
	"log"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/config"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/db"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/exchange"
	repository "github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/routes"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/services"
)

func main() {
	db.ConnectDatabase()

	cfg := config.LoadConfig()

	binanceConfig, err := cfg.GetExchangeConfig("binance")
	if err != nil {
		log.Fatalf("Failed to get Binance config: %v", err)
	}

	exchanges := map[string]exchange.Exchange{
		"binance": &exchange.Binance{
			Config: binanceConfig,
		},
	}
	walletBalanceRepo := repository.NewWalletBalanceRepository(db.DB)
	walletBalanceService := services.NewWalletBalanceService(walletBalanceRepo, exchanges)

	balanceHistoryRepo := repository.NewBalanceHistoryRepository(db.DB)
	balanceHistoryService := services.NewBalanceService(balanceHistoryRepo)

	router := routes.SetupRouter(walletBalanceService)

	cronService := services.NewCronService(walletBalanceService, balanceHistoryService)
	cronService.StartCronJobs()

	port := "8080"
	log.Printf("Starting server on port %s...", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
