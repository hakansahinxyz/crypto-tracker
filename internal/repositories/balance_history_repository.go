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

func (r *balanceHistoryRepository) CatchPumpDump() (models.Result, error) {

	query := `
    WITH TimeDifferences AS (
        SELECT 
            curr.id AS curr_id,
            curr.total_usd_value AS curr_value,
            curr.created_at AS curr_created_at,
            prev.id AS prev_id,
            prev.total_usd_value AS prev_value,
            prev.created_at AS prev_created_at,
            ABS(curr.total_usd_value - prev.total_usd_value) AS value_difference,
            (ABS(curr.total_usd_value - prev.total_usd_value) / prev.total_usd_value) * 100 AS percentage_difference,
            TIMESTAMPDIFF(MINUTE, prev.created_at, curr.created_at) AS time_difference
        FROM balance_histories curr
        JOIN balance_histories prev
          ON curr.created_at > prev.created_at
        WHERE curr.created_at >= NOW() - INTERVAL 15 MINUTE
					AND prev.created_at >= NOW() - INTERVAL 15 MINUTE
					-- AND TIMESTAMPDIFF(MINUTE, prev.created_at, curr.created_at) <= 15
          AND (ABS(curr.total_usd_value - prev.total_usd_value) / prev.total_usd_value) * 100 > 0.6
    )
    SELECT 
        prev_id,
        prev_value,
        prev_created_at,
        curr_id,
        curr_value,
        curr_created_at,
        value_difference,
        percentage_difference,
        time_difference
    FROM TimeDifferences
    ORDER BY value_difference DESC
    LIMIT 1;
    `

	var result models.Result
	if err := r.db.Raw(query).Scan(&result).Error; err != nil {
		panic(err)
	}

	return result, nil
}
