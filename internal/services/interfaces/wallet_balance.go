package services

import (
	"context"

	dto "github.com/hakansahinxyz/crypto-tracker-backend/internal/dto/request"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

type WalletBalanceService interface {
	GetAllBalances(ctx context.Context, filter dto.WalletBalanceRequest) ([]models.WalletBalance, error)
	UpdateBalance(ctx context.Context, balance *models.WalletBalance) error
	CreateBalance(ctx context.Context, balance *models.WalletBalance) error
}
