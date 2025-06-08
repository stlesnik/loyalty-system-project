package handler

import (
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"net/http"
)

type BalanceHandler struct {
	s *service.BalanceService
}

func NewBalanceHandler(s *service.BalanceService) *BalanceHandler { return &BalanceHandler{s: s} }

func (b *BalanceHandler) GetBalance(res http.ResponseWriter, req *http.Request)     {}
func (b *BalanceHandler) Withdraw(res http.ResponseWriter, req *http.Request)       {}
func (b *BalanceHandler) GetWithdrawals(res http.ResponseWriter, req *http.Request) {}
