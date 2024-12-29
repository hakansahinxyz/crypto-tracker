package dto

import (
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

type WalletBalanceRequest struct {
	AccountType  models.AccountType `json:"account_type"`
	ExchangeName string             `json:"exchange_name"`
	Asset        string             `json:"asset"`
}
