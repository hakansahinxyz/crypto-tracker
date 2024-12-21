package db

import (
	"log"

	"github.com/hakansahinxyz/crypto-tracker-backend/config"
	"github.com/hakansahinxyz/crypto-tracker-backend/models" // Burada models paketini ekliyoruz
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	cfg := config.LoadConfig()
	dsn := cfg.GetDBConnectionString()

	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = connection

	MigrateDatabase(DB)
}

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&models.Coin{}, &models.Exchange{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully.")
}
