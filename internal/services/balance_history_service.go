package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories/interfaces"
)

type BalanceService struct {
	balanceHistoryRepo interfaces.BalanceHistoryRepository
	tg                 *TelegramService
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

	if r.ValueDifference > 0 {
		var sb strings.Builder
		var isPump string
		if r.CurrValue > r.PrevValue {
			isPump = "+"
		} else {
			isPump = "-"
		}
		sb.WriteString(fmt.Sprintf("%.2f, %s%.2f, %s%.2f%%\n", r.CurrValue, isPump, r.ValueDifference, isPump, r.PercentageDifference))
		sb.WriteString(r.CurrCreatedAt.String())
		sb.WriteString("\n")
		sb.WriteString(r.PrevCreatedAt.String())
		s.tg.SendMessage(sb.String())
		log.Printf("%.2f  %.2f", r.ValueDifference, r.PercentageDifference)

		return r, nil
	}

	return models.Result{}, nil
}
