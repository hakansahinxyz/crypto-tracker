package services

import (
	"log"

	"github.com/robfig/cron/v3"
)

type CronService struct {
	walletBalanceService *WalletBalanceService
}

func NewCronService(walletBalanceService *WalletBalanceService) *CronService {
	return &CronService{
		walletBalanceService: walletBalanceService,
	}
}

func (cs *CronService) StartCronJobs() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	c := cron.New(cron.WithSeconds())

	for exchangeName := range cs.walletBalanceService.exchanges {
		exName := exchangeName

		c.AddFunc("10 */1 * * * *", func() {
			log.Printf("%s: Fetching Spot Wallet Balances...", exName)
			err := cs.walletBalanceService.FetchSpotWalletBalances(exName)
			if err != nil {
				log.Printf("Failed to fetch spot balances for %s: %v", exName, err)
			}
		})

		c.AddFunc("20 */1 6 * * *", func() {
			log.Printf("%s: Fetching Margin Wallet Balances...", exName)
			err := cs.walletBalanceService.FetchMarginWalletBalances(exName)
			if err != nil {
				log.Printf("Failed to fetch margin balances for %s: %v", exName, err)
			}
		})

		c.AddFunc("50 * 6 * * *", func() {
			log.Printf("%s: Fetching Future Wallet Balances...", exName)
			err := cs.walletBalanceService.FetchFutureAccountBalance(exName)
			if err != nil {
				log.Printf("Failed to fetch future balances for %s: %v", exName, err)
			}
		})

		c.AddFunc("0 * 6 * * *", func() {
			log.Printf("%s: Calculating Total Wallet Balances...", exName)
			err := cs.walletBalanceService.CalculateTotalUSDBalance(exName)
			if err != nil {
				log.Printf("Failed to calculate total balances for %s: %v", exName, err)
			}
		})
	}

	c.Start()
	select {}
}
