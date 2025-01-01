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

type AccountInfo struct {
	Balances []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
}

const (
	baseSpotURL   = "https://api.binance.com"
	baseMarginURL = "https://papi.binance.com"

	spotAccountInfoURL      = "%s/api/v3/account?%s"
	marginAccountBalanceURL = "%s/papi/v1/balance?%s"
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

	url := fmt.Sprintf(spotAccountInfoURL, baseSpotURL, queryStringWithSignature)

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
