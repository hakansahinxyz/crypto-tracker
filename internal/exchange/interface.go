package exchange

import "github.com/hakansahinxyz/crypto-tracker-backend/internal/models"

type Exchange interface {
	FetchSpotWalletBalances() ([]models.WalletBalance, error)
	FetchMarginWalletBalances() ([]models.WalletBalance, error)
	FetchFutureAccountBalance() ([]models.WalletBalance, error)
	CalculateTotalUSDBalance() (float64, error)
}
