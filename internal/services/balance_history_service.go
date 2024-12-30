package services

import (
	"context"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories/interfaces"
)

type BalanceService struct {
	balanceHistoryRepo interfaces.BalanceHistoryRepository
}

func NewBalanceService(balanceHistoryRepo interfaces.BalanceHistoryRepository) *BalanceService {
	return &BalanceService{balanceHistoryRepo: balanceHistoryRepo}
}

func (s *BalanceService) SaveCurrentBalance(ctx context.Context, usdValue float64) error {
	history := &models.BalanceHistory{
		ExchangeID:    1,
		TotalUSDValue: usdValue,
	}
	return s.balanceHistoryRepo.SaveBalanceHistory(ctx, history)
}

func (s *BalanceService) GetRecentHistory(ctx context.Context, limit int) ([]models.BalanceHistory, error) {
	return s.balanceHistoryRepo.GetBalanceHistory(ctx, limit)
}

func (s *BalanceService) GetActualBalance(ctx context.Context) (*models.BalanceHistory, error) {
	return s.balanceHistoryRepo.GetActualBalance(ctx)
}
