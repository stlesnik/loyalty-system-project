package handler

import (
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"net/http"
)

type OrderHandler struct {
	s *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler { return &OrderHandler{s: s} }

func (o *OrderHandler) UploadOrder(res http.ResponseWriter, req *http.Request) {}
func (o *OrderHandler) GetOrders(res http.ResponseWriter, req *http.Request)   {}
