// services/coinService.go
package services

import (
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/db"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

func CreateExchange(exchange *models.Exchange) error {
	return db.DB.Create(exchange).Error
}
