package exchange

import (
	"fmt"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

type BaseExchange struct{}

func (b *BaseExchange) FetchSpotWalletBalances() ([]models.WalletBalance, error) {
	return nil, fmt.Errorf("FetchSpotWalletBalances is not implemented")
}

func (b *BaseExchange) FetchMarginWalletBalances() ([]models.WalletBalance, error) {
	return nil, fmt.Errorf("FetchMarginWalletBalances is not implemented")
}

func (b *BaseExchange) FetchFutureAccountBalance() ([]models.WalletBalance, error) {
	return nil, fmt.Errorf("FetchFutureAccountBalance is not implemented")
}

func (b *BaseExchange) CalculateTotalUSDBalance(balances []models.WalletBalance) (float64, error) {
	return 0, fmt.Errorf("CalculateTotalUSDBalance is not implemented")
}
