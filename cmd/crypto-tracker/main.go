// main.go
package main

import (
	"log"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/db"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/exchange"
	repository "github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/routes"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/services"
)

func main() {
	db.ConnectDatabase()

	exchanges := map[string]exchange.Exchange{
		"binance": &exchange.Binance{
			ApiKey:    "WcgWNaKfrbthff5fJmpPG7SREvR0CPhq8Ucijthy7cfKwpgheab9RLzH1VfUpw5I",
			SecretKey: "bkmcU2z2CDPdGueOd8TG6N6lCjNsdlahQGlMpig9Z3SVqCXtQ1kabtKpWEy91J1h",
		},
	}
	walletBalanceRepo := repository.NewWalletBalanceRepository(db.DB)
	walletBalanceService := services.NewWalletBalanceService(walletBalanceRepo, exchanges)

	router := routes.SetupRouter(walletBalanceService)

	cronService := services.NewCronService(walletBalanceService)
	cronService.StartCronJobs()

	port := "8080"
	log.Printf("Starting server on port %s...", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
