package dto

import (
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

type WalletBalanceRequest struct {
	AccountType models.AccountType `form:"account_type" validate:"omitempty,oneof=spot margin futures"`
	Asset       string             `form:"asset"`
}
