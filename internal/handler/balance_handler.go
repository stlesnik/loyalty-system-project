package handler

import (
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"net/http"
)

type BalanceHandler struct {
	s   *service.BalanceService
	cfg *config.Config
}

func NewBalanceHandler(s *service.BalanceService, cfg *config.Config) *BalanceHandler {
	return &BalanceHandler{s: s, cfg: cfg}
}

func (b *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request)     {}
func (b *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request)       {}
func (b *BalanceHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {}
