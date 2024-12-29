package mappers

import (
	dto "github.com/hakansahinxyz/crypto-tracker-backend/internal/dto/resource"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

func ToBalanceHistoryResponse(balance models.BalanceHistory, exchangeName string) dto.BalanceHistoryResource {
	return dto.BalanceHistoryResource{
		ExchangeName: exchangeName,
		Timestamp:    balance.CreatedAt,
		BalanceUSD:   balance.TotalUSDValue,
	}
}

func ToWalletBalanceResponse(balance models.WalletBalance) dto.WalletBalanceResource {
	return dto.WalletBalanceResource{
		ID:          balance.ID,
		AccountType: balance.AccountType,
		Exchange:    balance.Exchange,
		Asset:       balance.Asset,
		Amount:      balance.Amount,
		USDValue:    balance.USDValue,
		IsActive:    balance.IsActive,
		UpdatedAt:   balance.UpdatedAt,
	}
}

func ToWalletBalancesResponse(balances []models.WalletBalance) dto.WalletBalancesResource {
	r := make([]dto.WalletBalanceResource, 0)

	for _, balance := range balances {
		r = append(r, ToWalletBalanceResponse(balance))
	}

	return r
}
