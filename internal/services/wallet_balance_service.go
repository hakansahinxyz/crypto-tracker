package services

import (
	"context"

	dto "github.com/hakansahinxyz/crypto-tracker-backend/internal/dto/request"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/repositories/interfaces"
)

type WalletBalanceService struct {
	repo interfaces.WalletBalanceRepository
}

func NewWalletBalanceService(repo interfaces.WalletBalanceRepository) *WalletBalanceService {
	return &WalletBalanceService{repo: repo}
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

func (s *WalletBalanceService) FetchSpotWalletBalances() {
	fetchSpotWalletBalancesFromBinance()
}

func (s *WalletBalanceService) FetchMarginWalletBalances() {
	fetchMarginWalletBalancesFromBinance()
}

func (s *WalletBalanceService) FetchFutureAccountBalance() {
	fetchFutureAccountBalance()
}

func (s *WalletBalanceService) CalculateTotalUSDBalance() {
	CalculateTotalUSDBalance()
}
