// config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

type ExchangeConfig struct {
	ApiKey    string
	SecretKey string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASS"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func (c *Config) GetExchangeConfig(exchangeName string) (*ExchangeConfig, error) {
	switch exchangeName {
	case "binance":
		return &ExchangeConfig{
			ApiKey:    os.Getenv("BINANCE_API_KEY"),
			SecretKey: os.Getenv("BINANCE_SECRET_KEY"),
		}, nil
	default:
		return nil, fmt.Errorf("exchange %s not supported", exchangeName)
	}
}
