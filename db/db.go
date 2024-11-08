// db/db.go
package db

import (
	"log"

	"github.com/hakansahinxyz/crypto-tracker-backend/config"
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
	log.Println("Database connection established")
}
