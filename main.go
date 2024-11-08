// main.go
package main

import (
	"github.com/hakansahinxyz/crypto-tracker-backend/db"
)

func main() {
	// MySQL'e bağlantı kuruyoruz
	db.ConnectDatabase()

	// Diğer yapılandırmalar ve API başlatma işlemleri buraya eklenecek
}
