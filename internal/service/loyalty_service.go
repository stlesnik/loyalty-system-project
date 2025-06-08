package service

import "github.com/stlesnik/loyalty-system-project/internal/repository"

type BalanceService struct {
	repBal  repository.BalanceRepository
	repWith repository.WithdrawalRepository
}

func NewBalanceService(repBal repository.BalanceRepository, repWith repository.WithdrawalRepository) *BalanceService {
	return &BalanceService{repBal: repBal, repWith: repWith}
}
