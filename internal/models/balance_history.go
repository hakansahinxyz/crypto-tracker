package models

import (
	"time"
)

type BalanceHistory struct {
	ID            uint      `gorm:"primaryKey"`
	ExchangeID    uint      `gorm:"not null"`
	TotalUSDValue float64   `gorm:"type:decimal(18,2);not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}
