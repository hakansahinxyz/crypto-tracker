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

	c.AddFunc("10 */5 * * * *", func() {
		log.Println("Fetching Spot Wallet Balances...")
		cs.walletBalanceService.FetchSpotWalletBalances()
	})

	c.AddFunc("20 */5 * * * *", func() {
		log.Println("Fetching Margin Wallet Balances...")
		cs.walletBalanceService.FetchMarginWalletBalances()
	})

	c.AddFunc("50 * * * * *", func() {
		log.Println("Fetching Future Wallet Balances...")
		cs.walletBalanceService.FetchFutureAccountBalance()
	})

	c.AddFunc("0 * * * * *", func() {
		log.Println("Calculating Total USD Balance...")
		cs.walletBalanceService.CalculateTotalUSDBalance()
	})

	c.Start()

	log.Println("Cron jobs started.")
	select {} // Prevent the main goroutine from exiting
}
