package handler

import (
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"net/http"
)

type OrderHandler struct {
	s   *service.OrderService
	cfg *config.Config
}

func NewOrderHandler(s *service.OrderService, cfg *config.Config) *OrderHandler {
	return &OrderHandler{s: s, cfg: cfg}
}

func (o *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {}
func (o *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request)   {}
