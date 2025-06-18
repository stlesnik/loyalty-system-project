package service

import (
	"context"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
)

type OrderService struct {
	rep OrderRepository
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

func NewOrderService(rep OrderRepository) *OrderService {
	return &OrderService{rep: rep}
}
