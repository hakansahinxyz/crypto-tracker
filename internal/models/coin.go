// models/coin.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Coin struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `json:"user_id"`     // Her kullanıcıya özel hesap
	Symbol     string         `json:"symbol"`      // Örneğin: BTC, ETH
	TotalBuy   float64        `json:"total_buy"`   // Toplam alım miktarı
	TotalSell  float64        `json:"total_sell"`  // Toplam satış miktarı
	ProfitLoss float64        `json:"profit_loss"` // Kâr veya zarar durumu
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete için
}
