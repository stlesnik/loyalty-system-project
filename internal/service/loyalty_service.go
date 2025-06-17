package service

type BalanceService struct {
	repBal  BalanceRepository
	repWith WithdrawalRepository
}

func NewBalanceService(repBal BalanceRepository, repWith WithdrawalRepository) *BalanceService {
	return &BalanceService{repBal: repBal, repWith: repWith}
}
