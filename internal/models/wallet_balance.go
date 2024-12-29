package models

import (
	"time"
)

type AccountType string

const (
	AccountTypeSpot    AccountType = "spot"
	AccountTypeMargin  AccountType = "margin"
	AccountTypeFutures AccountType = "futures"
)

type WalletBalance struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	AccountType AccountType `json:"account_type" gorm:"type:enum('spot', 'margin', 'futures');default:'spot';not null;index:unique_wallet_balance,unique"`
	ExchangeID  uint        `json:"exchange_id" gorm:"not null;index:unique_wallet_balance,unique"`
	Exchange    Exchange    `json:"exchange" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Asset       string      `json:"asset" gorm:"size:10;not null;index:unique_wallet_balance,unique"`
	Amount      float64     `json:"amount" gorm:"type:decimal(18,8);not null"`
	USDValue    float64     `json:"usd_value" gorm:"type:decimal(18,2)"`
	IsActive    bool        `json:"is_active" gorm:"default:true;not null"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}
