package repository

import (
	"context"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type balanceHistoryRepository struct {
	db *gorm.DB
}

// NewBalanceHistoryRepository creates a new BalanceHistoryRepository instance.
func NewBalanceHistoryRepository(db *gorm.DB) interfaces.BalanceHistoryRepository {
	return &balanceHistoryRepository{db: db}
}

func (r *balanceHistoryRepository) GetBalanceHistory(ctx context.Context, limit int) ([]models.BalanceHistory, error) {
	var history []models.BalanceHistory
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (r *balanceHistoryRepository) SaveBalanceHistory(ctx context.Context, history *models.BalanceHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *balanceHistoryRepository) DeleteOldHistory(ctx context.Context, threshold int) error {
	return r.db.WithContext(ctx).
		Where("id <= ?", threshold).
		Delete(&models.BalanceHistory{}).Error
}

func (r *balanceHistoryRepository) GetActualBalance(ctx context.Context) (*models.BalanceHistory, error) {
	var balance *models.BalanceHistory
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(1).
		Find(balance).Error; err != nil {
		return nil, err
	}
	return balance, nil
}
