package models

import (
	"time"
)

type BalanceHistory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ExchangeID    uint      `json:"exchange_id" gorm:"not null"`
	TotalUSDValue float64   `json:"total_usd_value" gorm:"type:decimal(18,2);not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}
