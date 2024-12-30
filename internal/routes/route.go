// routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/hakansahinxyz/crypto-tracker-backend/internal/handlers"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/services"
)

func SetupRouter(walletBalanceService *services.WalletBalanceService) *gin.Engine {
	r := gin.Default()

	walletBalanceHandlers := handlers.NewWalletBalanceHandler(walletBalanceService)

	walletBalancesRoutes := r.Group("/wallet-balance")
	{
		walletBalancesRoutes.GET("/", walletBalanceHandlers.GetAllBalances)
	}

	/* coinRoutes := r.Group("/coins")
	{
		coinRoutes.POST("/", controllers.CreateCoin)
		coinRoutes.GET("/", controllers.GetAllCoins)
		coinRoutes.GET("/:id", controllers.GetCoinByID)
		coinRoutes.PUT("/:id", controllers.UpdateCoin)
		coinRoutes.DELETE("/:id", controllers.DeleteCoin)
	} */

	/* exchangeRoutes := r.Group("/exchanges")
	{
		exchangeRoutes.POST("/", controllers.CreateExchange)
	} */

	// Transaction işlemleri için route grubu
	/* transactionRoutes := r.Group("/transactions")
	{
		transactionRoutes.POST("/", controllers.CreateTransaction)                   // Yeni bir transaction ekleme
		transactionRoutes.GET("/coin/:coin_id", controllers.GetTransactionsByCoinID) // Belirli bir coine ait tüm işlemleri listeleme
	} */

	return r
}
