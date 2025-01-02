package exchange

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/config"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
)

type Binance struct {
	BaseExchange
	Config *config.ExchangeConfig
}

const (
	baseURL = "https://api.binance.com"

	spotAccountInfoURL      = "%s/api/v3/account?%s"
	marginAccountBalanceURL = "%s/sapi/v1/margin/account?%s"
)

func (b *Binance) signQuery(query string) string {
	mac := hmac.New(sha256.New, []byte(b.Config.SecretKey))
	mac.Write([]byte(query))
	return hex.EncodeToString(mac.Sum(nil))
}

func (b *Binance) FetchSpotWalletBalances() ([]models.WalletBalance, error) {

	timestamp := time.Now().UnixMilli()
	queryString := fmt.Sprintf("timestamp=%d&omitZeroBalances=true", timestamp)
	signature := b.signQuery(queryString)
	queryStringWithSignature := fmt.Sprintf("%s&signature=%s", queryString, signature)

	url := fmt.Sprintf(spotAccountInfoURL, baseURL, queryStringWithSignature)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create Request: %v", err)
	}

	req.Header.Set("X-MBX-APIKEY", b.Config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to fetch account balance:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API Error: %s", body)
		return nil, err
	}

	var accountInfo AccountInfo
	if err := json.Unmarshal(body, &accountInfo); err != nil {
		log.Println("Failed to unmarshal account info:", err)
		return nil, err
	}

	var walletBalances []models.WalletBalance

	for _, balance := range accountInfo.Balances {

		free, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			continue
		}

		locked, err := strconv.ParseFloat(balance.Locked, 64)
		if err != nil {
			continue
		}

		amount := free + locked

		walletBalances = append(walletBalances, models.WalletBalance{
			ExchangeID:  1,
			Asset:       balance.Asset,
			Amount:      amount,
			AccountType: models.AccountTypeSpot,
		})
	}

	return walletBalances, nil
}

func (b *Binance) FetchMarginWalletBalances() ([]models.WalletBalance, error) {
	timestamp := time.Now().UnixMilli()
	queryString := fmt.Sprintf("timestamp=%d&omitZeroBalances=true", timestamp)
	signature := b.signQuery(queryString)
	queryStringWithSignature := fmt.Sprintf("%s&signature=%s", queryString, signature)

	url := fmt.Sprintf(marginAccountBalanceURL, baseURL, queryStringWithSignature)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create Request: %v", err)
	}

	req.Header.Set("X-MBX-APIKEY", b.Config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to fetch margin account balance:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API Error: %s", body)
		return nil, err
	}

	var account MarginAccountResponse
	if err := json.Unmarshal(body, &account); err != nil {
		log.Println("Failed to unmarshal margin balances:", err)
		return nil, err
	}

	var walletBalances []models.WalletBalance

	for _, balance := range account.Assets {
		amount, err := strconv.ParseFloat(balance.NetAsset, 64)
		if err != nil {
			continue
		}
		if amount > 0 || balance.Asset == "USDT" {
			walletBalances = append(walletBalances, models.WalletBalance{
				ExchangeID:  1,
				Asset:       balance.Asset,
				Amount:      amount,
				AccountType: models.AccountTypeMargin,
			})
		}

	}

	return walletBalances, nil
}

func (b *Binance) FetchFutureAccountBalance() ([]models.WalletBalance, error) {
	timestamp := time.Now().UnixMilli()
	queryString := fmt.Sprintf("timestamp=%d&recvWindow=%d", timestamp, 7000)
	signature := b.signQuery(queryString)
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v3/balance?%s&signature=%s", queryString, signature)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create Request: %v", err)
	}

	req.Header.Set("X-MBX-APIKEY", b.Config.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to fetch account balance:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API Error: %s", body)
		return nil, err
	}

	var balances FutureBalanceResponse
	if err := json.Unmarshal(body, &balances); err != nil {
		log.Println("Failed to unmarshal account info:", err)
		return nil, err
	}

	var walletBalances []models.WalletBalance
	for _, balance := range balances {
		if balance.Asset == "USDT" {
			value, err := strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				continue
			}
			pnl, err := strconv.ParseFloat(balance.UnrealizedPNL, 64)
			if err != nil {
				continue
			}
			log.Printf("Future Balance: %.2f USD", value)
			log.Printf("Future UnPNL  : %.2f USD", pnl)
			log.Printf("Future Total  : %.2f USD, UnPNL  : %.2f USD", value+pnl, pnl)

			walletBalances = append(walletBalances, models.WalletBalance{
				ExchangeID:  1,
				Asset:       balance.Asset,
				Amount:      value + pnl,
				AccountType: models.AccountTypeFutures,
			})
		}
	}

	return walletBalances, nil
}

func (b *Binance) CalculateTotalUSDBalance(balances []models.WalletBalance) (float64, error) {

	priceMap, err := FetchPrices()
	if err != nil {
		log.Printf("error fetching prices: %v\n", err)
		return 0, err
	}

	totalBalance := 0.0
	for _, balance := range balances {
		if balance.Asset == "USDT" {
			totalBalance += balance.Amount
			continue
		}
		symbol := balance.Asset + "USDT"
		price, ok := priceMap[symbol]
		if !ok {
			log.Printf("Price not found for asset: %s", balance.Asset)
			continue
		}
		totalBalance += balance.Amount * price
	}
	log.Printf("Total Balance: %.2f", totalBalance)

	return totalBalance, nil
}

func FetchPrices() (map[string]float64, error) {
	url := "https://api.binance.com/api/v3/ticker/price"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching prices: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var prices []Price
	if err := json.Unmarshal(body, &prices); err != nil {
		return nil, fmt.Errorf("error unmarshaling prices: %v", err)
	}

	priceMap := make(map[string]float64)
	for _, price := range prices {
		priceMap[price.Symbol] = price.Price
	}
	return priceMap, nil
}
