package service

import (
	"context"
	"github.com/stlesnik/loyalty-system-project/internal/model"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
)

type OrderService struct {
	rep OrderRepository
}

func NewOrderService(rep OrderRepository) *OrderService {
	return &OrderService{rep: rep}
}

func (s OrderService) UploadOrder(ctx context.Context, userID int, orderNumber string) error {
	existing, err := s.rep.GetByID(ctx, orderNumber)

	if err == nil {
		if existing.UserID == userID {
			return utils.ErrOrderAlreadyUploaded
		}
		return utils.ErrOrderConflict
	}

	_, err = s.rep.Create(ctx, userID, orderNumber)
	return err
}

func (s OrderService) GetUserOrders(ctx context.Context, userID int) ([]model.Order, error) {
	orders, err := s.rep.GetByUserID(ctx, userID)
	switch {
	case len(orders) == 0:
		return nil, utils.ErrNoOrders
	case err != nil:
		return nil, err
	default:
		return orders, nil
	}
}
