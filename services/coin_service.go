// services/coinService.go
package services

import (
	"github.com/hakansahinxyz/crypto-tracker-backend/db"
	"github.com/hakansahinxyz/crypto-tracker-backend/models"
)

func CreateCoin(coin *models.Coin) error {
	return db.DB.Create(coin).Error
}

func GetAllCoins() ([]models.Coin, error) {
	var coins []models.Coin
	err := db.DB.Find(&coins).Error
	return coins, err
}

func GetCoinByID(id uint) (models.Coin, error) {
	var coin models.Coin
	err := db.DB.First(&coin, id).Error
	return coin, err
}

func UpdateCoin(coin *models.Coin) error {
	return db.DB.Save(coin).Error
}

func DeleteCoin(id uint) error {
	return db.DB.Delete(&models.Coin{}, id).Error
}
