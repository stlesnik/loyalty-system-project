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
			utils.Log.Infow("Order already exists for this user", "order_number", orderNumber, "db_user_id", existing.UserID, "req_user_id", userID)
			return utils.ErrOrderAlreadyUploaded
		}
		utils.Log.Infow("Order already exists for different user", "order_number", orderNumber, "db_user_id", existing.UserID, "req_user_id", userID)
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
