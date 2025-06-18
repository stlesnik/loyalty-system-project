package service

import (
	"context"
	"github.com/stlesnik/loyalty-system-project/internal/client"
	"github.com/stlesnik/loyalty-system-project/internal/model"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"go.uber.org/zap"
	"time"
)

const FetchAccrualTimeout = 30 * time.Second

type OrderService struct {
	rep        OrderRepository
	accrualCli *client.AccrualClient
	workerPool chan string
}

func NewOrderService(rep OrderRepository, accrualClient *client.AccrualClient) *OrderService {
	return &OrderService{
		rep:        rep,
		accrualCli: accrualClient,
		workerPool: make(chan string, 1000),
	}
}

func (s *OrderService) StartWorkers(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go s.accrualWorker(ctx, i)
	}
	utils.Log.Infow("Accrual workers started", "count", workerCount)
}

func (s *OrderService) accrualWorker(ctx context.Context, id int) {
	utils.Log.Debugw("Worker started", "id", id)

	for {
		select {
		case orderNumber := <-s.workerPool:
			s.processOrder(ctx, orderNumber)
		case <-ctx.Done():
			utils.Log.Debugw("Worker stopped", "id", id)
			return
		}
	}
}

func (s *OrderService) processOrder(ctx context.Context, orderNumber string) {
	ctx, cancel := context.WithTimeout(ctx, FetchAccrualTimeout)
	defer cancel()

	resp, err := s.accrualCli.FetchAccrual(ctx, orderNumber)
	if err != nil {
		utils.Log.Warnw("Accrual fetch failed", "order", orderNumber, "err", err)
		return
	}

	status := resp.Status
	if resp.Status == "REGISTERED" {
		status = "NEW"
	}
	if resp.Status == "PROCESSED" && resp.Accrual == nil {
		utils.Log.Errorw("Invalid PROCESSED without accrual", "order", orderNumber)
		status = "INVALID"
	}
	if err := s.rep.Update(ctx, orderNumber, status, resp.Accrual); err != nil {
		utils.Log.Error("Order status update failed",
			zap.String("order", orderNumber),
			zap.Error(err))
	}

	if resp.Status == "PROCESSING" || resp.Status == "REGISTERED" {
		time.AfterFunc(30*time.Second, func() {
			s.workerPool <- orderNumber
		})
	}
}

func (s *OrderService) UploadOrder(ctx context.Context, userID int, orderNumber string) error {
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
	s.workerPool <- orderNumber
	return err
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID int) ([]model.Order, error) {
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
