package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type BinanceWebSocketResponse struct {
	Symbol string `json:"s"`
	Price  string `json:"p"`
}

func StartPriceService() {
	// Binance WebSocket URL
	url := "wss://stream.binance.com:9443/ws/btcusdt@trade"

	// WebSocket bağlantısı oluştur
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket:", err)
	}
	defer conn.Close()

	for {
		// Mesajları dinle
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		// Gelen mesajı çözümle
		var response BinanceWebSocketResponse
		if err := json.Unmarshal(message, &response); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		// Bitcoin fiyatını yazdır
		fmt.Printf("Current Bitcoin Price: %s USDT\n", response.Price)
	}
}
