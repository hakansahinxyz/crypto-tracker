package dto

import (
	"time"
)

type BalanceHistoryResource struct {
	ExchangeName string    `json:"exchange_name"`
	Timestamp    time.Time `json:"timestamp"`
	BalanceUSD   float64   `json:"balance_usd"`
}
