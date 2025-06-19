package handler

import (
	"encoding/json"
	"errors"
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/middleware"
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"io"
	"net/http"
	"strings"
)

type OrderHandler struct {
	s   service.OrderService
	cfg *config.Config
}

func NewOrderHandler(s service.OrderService, cfg *config.Config) *OrderHandler {
	return &OrderHandler{s: s, cfg: cfg}
}

func (oH *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	orderNumber := strings.TrimSpace(string(body))

	if !utils.ValidLuhn(orderNumber) {
		http.Error(w, "Invalid order format", http.StatusUnprocessableEntity)
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKeyName).(int)
	err = oH.s.UploadOrder(r.Context(), userID, orderNumber)
	switch {
	case errors.Is(err, utils.ErrOrderAlreadyUploaded):
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, utils.ErrOrderConflict):
		http.Error(w, "Order conflict", http.StatusConflict)
	case err != nil:
		http.Error(w, "Server error", http.StatusInternalServerError)
	default:
		utils.Log.Infow("Order uploaded", "order", orderNumber, "user_id", userID)
		w.WriteHeader(http.StatusAccepted)
	}
}

func (oH *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKeyName).(int)
	orders, err := oH.s.GetUserOrders(r.Context(), userID)
	switch {
	case errors.Is(err, utils.ErrNoOrders):
		w.WriteHeader(http.StatusNoContent)
		return
	case err != nil:
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(orders)
		if err != nil {
			utils.Log.Errorf("json encode error: %s", err.Error())
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	}
}
