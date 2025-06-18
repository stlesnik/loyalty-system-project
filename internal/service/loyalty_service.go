package service

import (
	"context"
	"fmt"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
)

type BalanceService struct {
	repBal  BalanceRepository
	repWith WithdrawalRepository
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func NewBalanceService(repBal BalanceRepository, repWith WithdrawalRepository) *BalanceService {
	return &BalanceService{repBal: repBal, repWith: repWith}
}

func (s *BalanceService) GetBalance(ctx context.Context, userID int) (*Balance, error) {
	current, err := s.repBal.GetTotal(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get current failed: %w", err)
	}

	withdrawn, err := s.repWith.GetTotal(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get withdrawn failed: %w", err)
	}

	return &Balance{
		Current:   current,
		Withdrawn: withdrawn,
	}, nil
}

func (s *BalanceService) CreateWithdrawal(ctx context.Context, userID int, orderNumber string, amount float64) error {
	current, err := s.repBal.GetTotal(ctx, userID)
	if err != nil {
		return fmt.Errorf("get current failed: %w", err)
	}
	if current < amount {
		return utils.ErrInsufficientFunds
	}
	_, err = s.repWith.Create(ctx, userID, orderNumber, amount)
	if err != nil {
		return err
	}
	return nil
}
