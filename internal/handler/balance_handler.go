package handler

import (
	"encoding/json"
	"errors"
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/middleware"
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"net/http"
)

type BalanceHandler struct {
	s   service.BalanceService
	cfg *config.Config
}

func NewBalanceHandler(s service.BalanceService, cfg *config.Config) *BalanceHandler {
	return &BalanceHandler{s: s, cfg: cfg}
}

func (bH *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKeyName).(int)
	balance, err := bH.s.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	utils.Log.Infow("Got balance", "user_id", userID, "balance", balance)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(balance)
	if err != nil {
		utils.Log.Errorf("json encode error: %s", err.Error())
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func (bH *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKeyName).(int)
	var req struct {
		Order  string  `json:"order"`
		Amount float64 `json:"sum"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if !utils.ValidLuhn(req.Order) {
		http.Error(w, "Invalid order format", http.StatusUnprocessableEntity)
		return
	}

	err := bH.s.CreateWithdrawal(r.Context(), userID, req.Order, req.Amount)
	switch {
	case errors.Is(err, utils.ErrInsufficientFunds):
		http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
	case errors.Is(err, utils.ErrOrderAlreadyUploaded):
		http.Error(w, "Order conflict", http.StatusUnprocessableEntity)
	case err != nil:
		utils.Log.Errorf("create withdrawal error: %s", err.Error())
		http.Error(w, "Server error", http.StatusInternalServerError)
	default:
		utils.Log.Infow("Created withdrawal", "user_id", userID, "amount", req.Amount)
		w.WriteHeader(http.StatusOK)
	}
}

func (bH *BalanceHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKeyName).(int)
	withdrawals, err := bH.s.GetAllWithdrawals(r.Context(), userID)
	if err != nil {
		utils.Log.Errorf("get withdrawals error: %s", err.Error())
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(withdrawals)
	if err != nil {
		utils.Log.Errorf("json encode error: %s", err.Error())
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
