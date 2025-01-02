package repository

import (
	"context"

	dto "github.com/hakansahinxyz/crypto-tracker-backend/internal/dto/request"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type walletBalanceRepository struct {
	db *gorm.DB
}

func NewWalletBalanceRepository(db *gorm.DB) interfaces.WalletBalanceRepository {
	return &walletBalanceRepository{db: db}
}

func (r *walletBalanceRepository) GetAllBalances(ctx context.Context, req dto.WalletBalanceRequest) ([]models.WalletBalance, error) {
	var balances []models.WalletBalance
	query := r.db.WithContext(ctx).Model(&models.WalletBalance{})

	if req.AccountType != "" {
		query = query.Where("account_type = ?", req.AccountType)
	}

	if req.Asset != "" {
		query = query.Where("asset = ?", req.Asset)
	}

	if err := query.Find(&balances).Error; err != nil {
		return nil, err
	}
	return balances, nil
}

func (r *walletBalanceRepository) GetActiveAndNonZeroBalances() ([]models.WalletBalance, error) {
	var balances []models.WalletBalance
	if err := r.db.Where("amount != 0 AND is_active = true").Find(&balances).Error; err != nil {
		return nil, err
	}
	return balances, nil
}
