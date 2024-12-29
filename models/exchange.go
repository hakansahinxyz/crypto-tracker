package models

import (
	"time"

	"gorm.io/gorm"
)

type Exchange struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `json:"name"`
	Balances  []WalletBalance `json:"balances" gorm:"foreignKey:ExchangeID"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}
