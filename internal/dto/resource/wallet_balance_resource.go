package dto

import (
	"time"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

type WalletBalanceResource struct {
	ID          uint               `json:"id"`
	AccountType models.AccountType `json:"account_type"`
	Exchange    models.Exchange    `json:"exchange"`
	Asset       string             `json:"asset"`
	Amount      float64            `json:"amount"`
	USDValue    float64            `json:"usd_value"`
	IsActive    bool               `json:"is_active"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type WalletBalancesResource []WalletBalanceResource
