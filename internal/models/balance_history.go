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

type Result struct {
	PrevID               uint      `gorm:"column:prev_id"`
	PrevValue            float64   `gorm:"column:prev_value"`
	PrevCreatedAt        time.Time `gorm:"column:prev_created_at"`
	CurrID               uint      `gorm:"column:curr_id"`
	CurrValue            float64   `gorm:"column:curr_value"`
	CurrCreatedAt        time.Time `gorm:"column:curr_created_at"`
	ValueDifference      float64   `gorm:"column:value_difference"`
	PercentageDifference float64   `gorm:"column:percentage_difference"`
	TimeDifference       int       `gorm:"column:time_difference"`
}
