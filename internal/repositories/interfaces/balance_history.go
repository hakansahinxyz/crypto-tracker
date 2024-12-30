package interfaces

import (
	"context"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

type BalanceHistoryRepository interface {
	GetBalanceHistory(ctx context.Context, limit int) ([]models.BalanceHistory, error)
	SaveBalanceHistory(ctx context.Context, history *models.BalanceHistory) error
	DeleteOldHistory(ctx context.Context, threshold int) error
	GetActualBalance(ctx context.Context) (*models.BalanceHistory, error)
}
