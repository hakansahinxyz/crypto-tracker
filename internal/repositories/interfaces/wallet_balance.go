package interfaces

import (
	"context"

	dto "github.com/hakansahinxyz/crypto-tracker-backend/internal/dto/request"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

type WalletBalanceRepository interface {
	GetAllBalances(ctx context.Context, req dto.WalletBalanceRequest) ([]models.WalletBalance, error)
}
