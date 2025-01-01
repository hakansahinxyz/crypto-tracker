package services

import (
	"context"
	"fmt"

	dto "github.com/hakansahinxyz/crypto-tracker-backend/internal/dto/request"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/exchange"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories/interfaces"
)

type WalletBalanceService struct {
	repo      interfaces.WalletBalanceRepository
	exchanges map[string]exchange.Exchange
}

func NewWalletBalanceService(
	repo interfaces.WalletBalanceRepository,
	exchanges map[string]exchange.Exchange,
) *WalletBalanceService {
	return &WalletBalanceService{repo: repo, exchanges: exchanges}
}

func (s *WalletBalanceService) GetAllBalances(ctx context.Context, req dto.WalletBalanceRequest) ([]models.WalletBalance, error) {
	return s.repo.GetAllBalances(ctx, req)
}

/* func (s *walletBalanceService) UpdateBalance(ctx context.Context, balance *models.WalletBalance) error {
    return s.repo.UpdateBalance(ctx, balance)
}

func (s *walletBalanceService) CreateBalance(ctx context.Context, balance *models.WalletBalance) error {
    return s.repo.CreateBalance(ctx, balance)
} */

func (s *WalletBalanceService) FetchSpotWalletBalances(exchangeName string) error {
	ex, exists := s.exchanges[exchangeName]
	if !exists {
		return fmt.Errorf("exchange %s not found", exchangeName)
	}
	balances, err := ex.FetchSpotWalletBalances()
	updateWalletBalances(models.AccountTypeSpot, balances)
	return err
}

func (s *WalletBalanceService) FetchMarginWalletBalances(exchangeName string) error {
	ex, exists := s.exchanges[exchangeName]
	if !exists {
		return fmt.Errorf("exchange %s not found", exchangeName)
	}
	_, err := ex.FetchMarginWalletBalances()
	return err
}

func (s *WalletBalanceService) FetchFutureAccountBalance(exchangeName string) error {
	ex, exists := s.exchanges[exchangeName]
	if !exists {
		return fmt.Errorf("exchange %s not found", exchangeName)
	}
	_, err := ex.FetchFutureAccountBalance()
	return err
}

func (s *WalletBalanceService) CalculateTotalUSDBalance(exchangeName string) error {
	ex, exists := s.exchanges[exchangeName]
	if !exists {
		return fmt.Errorf("exchange %s not found", exchangeName)
	}
	_, err := ex.CalculateTotalUSDBalance()
	return err
}
