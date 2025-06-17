package service

type OrderService struct {
	rep OrderRepository
}

func NewOrderService(rep OrderRepository) *OrderService {
	return &OrderService{rep: rep}
}
