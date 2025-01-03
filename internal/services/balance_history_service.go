package services

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories/interfaces"
)

type BalanceService struct {
	balanceHistoryRepo interfaces.BalanceHistoryRepository
	tg                 *TelegramService
	lastPumpDumb       models.Result
}

func NewBalanceService(
	balanceHistoryRepo interfaces.BalanceHistoryRepository,
	tg *TelegramService,
) *BalanceService {
	return &BalanceService{balanceHistoryRepo: balanceHistoryRepo, tg: tg}
}

func (s *BalanceService) SaveCurrentBalance(ctx context.Context, usdValue float64) error {
	history := &models.BalanceHistory{
		ExchangeID:    1,
		TotalUSDValue: usdValue,
	}
	return s.balanceHistoryRepo.SaveBalanceHistory(ctx, history)
}

func (s *BalanceService) GetRecentHistory(ctx context.Context, limit int) ([]models.BalanceHistory, error) {
	return s.balanceHistoryRepo.GetBalanceHistory(ctx, limit)
}

func (s *BalanceService) GetActualBalance(ctx context.Context) (*models.BalanceHistory, error) {
	return s.balanceHistoryRepo.GetActualBalance(ctx)
}

func (s *BalanceService) CatchPumpDump() (models.Result, error) {
	r, err := s.balanceHistoryRepo.CatchPumpDump()
	if err != nil {
		log.Printf("Failed to calculate PumpDump: %v", err)
		return models.Result{}, err
	}

	if reflect.DeepEqual(r, s.lastPumpDumb) {
		return models.Result{}, nil
	}

	if r.ValueDifference > 0 {
		direction := "-"
		status := "Dump"
		if r.CurrValue > r.PrevValue {
			direction = "+"
			status = "Pump"
		}

		message := fmt.Sprintf("🔥 Portfolio %s 🔥", status) +
			"\n" + fmt.Sprintf("Balance Change: %s$%.2f, %%%.2f", direction, r.ValueDifference, r.PercentageDifference) +
			"\n" + fmt.Sprintf("$%.2f -> $%.2f", r.PrevValue, r.CurrValue) +
			"\n" + fmt.Sprintf("%s -> %s", r.PrevCreatedAt.Format("15:04"), r.CurrCreatedAt.Format("15:04"))

		s.tg.SendMessage(message)
		log.Printf("%.2f  %.2f", r.ValueDifference, r.PercentageDifference)
		s.lastPumpDumb = r
		return r, nil
	}

	return models.Result{}, nil
}
