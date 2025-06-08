package service

import "github.com/stlesnik/loyalty-system-project/internal/repository"

type OrderService struct {
	rep repository.OrderRepository
}

func NewOrderService(rep repository.OrderRepository) *OrderService {
	return &OrderService{rep: rep}
}
