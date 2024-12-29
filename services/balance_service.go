package services

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

	"github.com/hakansahinxyz/crypto-tracker-backend/db"
	"github.com/hakansahinxyz/crypto-tracker-backend/models"
	cron "github.com/robfig/cron/v3"
)

const (
	apiKey    = "WcgWNaKfrbthff5fJmpPG7SREvR0CPhq8Ucijthy7cfKwpgheab9RLzH1VfUpw5I"
	secretKey = "bkmcU2z2CDPdGueOd8TG6N6lCjNsdlahQGlMpig9Z3SVqCXtQ1kabtKpWEy91J1h"
)

const (
	baseURL                 = "https://api.binance.com"
	baseMarginURL           = "https://papi.binance.com"
	spotAccountInfoURL      = "%s/api/v3/account?%s"
	marginAccountBalanceURL = "%s/papi/v1/balance?%s"
)

type AccountInfo struct {
	Balances []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
}

type CoinPrice struct {
	Price string `json:"price"`
}

type FutureBalance struct {
	Asset         string `json:"asset"`
	Balance       string `json:"balance"`
	UnrealizedPNL string `json:"crossUnPnl"`
}

type FutureBalanceResponse []FutureBalance

type MarginBalance struct {
	Asset    string `json:"asset"`
	Borrowed string `json:"borrowed"`
	Free     string `json:"free"`
	Interest string `json:"interest"`
	Locked   string `json:"locked"`
	NetAsset string `json:"netAsset"`
}

type MarginAccountResponse struct {
	MarginLevel                string          `json:"marginLevel"`
	CollateralMarginLevel      string          `json:"collateralMarginLevel"`
	TotalCollateralValueInUSDT string          `json:"totalCollateralValueInUSDT"`
	Assets                     []MarginBalance `json:"userAssets"`
}

type Price struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

func sign(queryString string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(queryString))
	return hex.EncodeToString(mac.Sum(nil))
}

func fetchSpotWalletBalancesFromBinance() {
	timestamp := time.Now().UnixMilli()
	queryString := fmt.Sprintf("timestamp=%d&omitZeroBalances=true", timestamp)
	signature := sign(queryString)
	queryStringWithSignature := fmt.Sprintf("%s&signature=%s", queryString, signature)

	url := fmt.Sprintf(spotAccountInfoURL, baseURL, queryStringWithSignature)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create Request: %v", err)
	}

	req.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to fetch account balance:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API Error: %s", body)
	}

	var accountInfo AccountInfo
	if err := json.Unmarshal(body, &accountInfo); err != nil {
		log.Println("Failed to unmarshal account info:", err)
		return
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

	updateWalletBalances(models.AccountTypeSpot, walletBalances)

	/* if err := db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "exchange_id"}, {Name: "asset"}},
		DoUpdates: clause.AssignmentColumns([]string{"amount", "updated_at"}),
	}).Create(&walletBalances).Error; err != nil {
		log.Fatalf("Failed to batch upsert wallet balances: %v", err)
	} */
}

func updateWalletBalances(accountType models.AccountType, walletBalances []models.WalletBalance) {
	updatedAssets := make(map[string]struct{})

	for _, balance := range walletBalances {
		result := db.DB.
			Where("exchange_id = ?", balance.ExchangeID).
			Where("asset = ?", balance.Asset).
			Where("account_type = ?", balance.AccountType).
			Assign(balance).
			FirstOrCreate(&balance)
		if result.Error != nil {
			log.Printf("Failed to save wallet balance for %s: %v", balance.Asset, result.Error)
		}

		if result.RowsAffected > 0 {
			updatedAssets[balance.Asset] = struct{}{}
		}

	}

	err := db.DB.Model(&models.WalletBalance{}).
		Where("exchange_id = ?", 1).
		Where("account_type = ?", accountType).
		Where("asset NOT IN ?", getAssetKeys(updatedAssets)).
		Where("is_active = true").
		Where("amount != 0").
		Update("amount", 0).Error
	if err != nil {
		log.Printf("Failed to reset balances for missing assets: %v", err)
	}

	//log.Printf("Successfully update %s balances ", accountType)
}

func getAssetKeys(assetSet map[string]struct{}) []string {
	keys := make([]string, 0, len(assetSet))
	for key := range assetSet {
		keys = append(keys, key)
	}
	return keys
}

func fetchMarginWalletBalancesFromBinance() {
	timestamp := time.Now().UnixMilli()
	//queryString := fmt.Sprintf("timestamp=%d", timestamp)
	queryString := fmt.Sprintf("timestamp=%d&omitZeroBalances=true", timestamp)
	signature := sign(queryString)
	//queryStringWithSignature := fmt.Sprintf("%s&signature=%s", queryString, signature)

	//url := fmt.Sprintf(marginAccountBalanceURL, baseMarginURL, queryStringWithSignature)
	url := fmt.Sprintf("https://api.binance.com/sapi/v1/margin/account?%s&signature=%s", queryString, signature)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create Request: %v", err)
	}

	req.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to fetch margin account balance:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API Error: %s", body)
	}

	var account MarginAccountResponse
	if err := json.Unmarshal(body, &account); err != nil {
		log.Println("Failed to unmarshal margin balances:", err)
		return
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

	updateWalletBalances(models.AccountTypeMargin, walletBalances)
}

func fetchAccountBalance() {

	timestamp := time.Now().UnixMilli()
	queryString := fmt.Sprintf("timestamp=%d&omitZeroBalances=true", timestamp)
	signature := sign(queryString)
	url := fmt.Sprintf("https://api.binance.com/api/v3/account?%s&signature=%s", queryString, signature)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create Request: %v", err)
	}

	req.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to fetch account balance:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API Error: %s", body)
	}

	var accountInfo AccountInfo
	if err := json.Unmarshal(body, &accountInfo); err != nil {
		log.Println("Failed to unmarshal account info:", err)
		return
	}

	totalUSDValue := 0.0
	for _, balance := range accountInfo.Balances {
		//fmt.Printf("%s: %s %s\n", balance.Asset, balance.Free, balance.Locked)
		free, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			continue
		}

		locked, err := strconv.ParseFloat(balance.Locked, 64)
		if err != nil {
			continue
		}

		amount := free + locked

		if amount == 0 {
			continue
		}

		if balance.Asset == "USDT" {
			totalUSDValue += amount
			continue
		}

		symbol := fmt.Sprintf("%sUSDT", balance.Asset)
		price, err := fetchCoinPrice(symbol)
		if err != nil {
			//log.Printf("%s failed to fetch price: %v", balance.Asset, err)
			continue
		}

		totalUSDValue += amount * price
	}

	log.Printf("Total USD Value: %.2f USD", totalUSDValue)
}

func fetchCoinPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch coin price: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read resp.body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API Error: %s", body)
	}

	var coinPrice CoinPrice
	if err := json.Unmarshal(body, &coinPrice); err != nil {
		return 0, fmt.Errorf("JSON parse error: %v", err)
	}

	price, err := strconv.ParseFloat(coinPrice.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("price parse error: %v", err)
	}

	return price, nil
}

func fetchFutureAccountBalance() {
	timestamp := time.Now().UnixMilli()
	queryString := fmt.Sprintf("timestamp=%d&recvWindow=%d", timestamp, 7000)
	signature := sign(queryString)
	url := fmt.Sprintf("https://fapi.binance.com/fapi/v3/balance?%s&signature=%s", queryString, signature)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create Request: %v", err)
	}

	req.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to fetch account balance:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API Error: %s", body)
	}

	var balances FutureBalanceResponse
	if err := json.Unmarshal(body, &balances); err != nil {
		log.Println("Failed to unmarshal account info:", err)
		return
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
			//log.Printf("Future Balance: %.2f USD", value)
			//log.Printf("Future UnPNL  : %.2f USD", pnl)
			//log.Printf("Future Total  : %.2f USD, UnPNL  : %.2f USD", value+pnl, pnl)

			walletBalances = append(walletBalances, models.WalletBalance{
				ExchangeID:  1,
				Asset:       balance.Asset,
				Amount:      value + pnl,
				AccountType: models.AccountTypeFutures,
			})
		}
	}

	updateWalletBalances(models.AccountTypeFutures, walletBalances)
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

func CalculateTotalUSDBalance() {
	var balances []models.WalletBalance
	if err := db.DB.Where("amount != 0 AND is_active = true").Find(&balances).Error; err != nil {
		log.Printf("error fetching balances: %v\n", err)
	}

	priceMap, err := FetchPrices()
	if err != nil {
		log.Printf("error fetching prices: %v\n", err)
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

	if err := SaveTotalBalance(totalBalance); err != nil {
		log.Printf("Error saving total balance: %v", err)
	}
	log.Printf("Total Balance: %.2f", totalBalance)

}

func SaveTotalBalance(totalBalance float64) error {
	history := models.BalanceHistory{
		ExchangeID:    1,
		TotalUSDValue: totalBalance,
	}

	if err := db.DB.Create(&history).Error; err != nil {
		return fmt.Errorf("error saving balance history: %v", err)
	}
	return nil
}

func StartBalanceService() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	c := cron.New(cron.WithSeconds())
	//c.AddFunc("@every 1m", fetchAccountBalance)
	//c.AddFunc("@every 3s", fetchFutureAccountBalance)
	//c.AddFunc("0 * * * * *", fetchAccountBalance)          // every minute; 12:00:00, 12:01:00, 12:02:00
	c.AddFunc("10 */5 * * * *", fetchSpotWalletBalancesFromBinance)   // every minute; 12:00:10, 12:01:10, 12:02:10
	c.AddFunc("20 */5 * * * *", fetchMarginWalletBalancesFromBinance) // every minute; 12:00:20, 12:01:20 12:02:20
	c.AddFunc("50 * * * * *", fetchFutureAccountBalance)              // every minute; 12:00:20, 12:01:20 12:02:20
	c.AddFunc("0 * * * * *", CalculateTotalUSDBalance)                // every minute; 12:00:20, 12:01:20 12:02:20

	c.Start()
}
