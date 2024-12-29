package mappers

import (
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/dto"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

func ToBalanceHistoryDTO(balance models.BalanceHistory, exchangeName string) dto.BalanceHistoryDTO {
	return dto.BalanceHistoryDTO{
		ExchangeName: exchangeName,
		Timestamp:    balance.CreatedAt,
		BalanceUSD:   balance.TotalUSDValue,
	}
}
