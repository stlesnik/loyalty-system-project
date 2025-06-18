package service

import (
	"context"
	"github.com/stlesnik/loyalty-system-project/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, login, password string) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByID(ctx context.Context, userID int) (*model.User, error)
}

type OrderRepository interface {
	Create(ctx context.Context, userID int, orderID string) (*model.Order, error)
	GetByUserID(ctx context.Context, userID int) ([]model.Order, error)
	GetByID(ctx context.Context, orderID string) (*model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

type BalanceRepository interface {
	Get(ctx context.Context, userID int) (*model.Balance, error)
	Update(ctx context.Context, userID int, deltaCurrent int, deltaWithdrawn int) error
}

type WithdrawalRepository interface {
	Create(ctx context.Context, userID int, orderID string, sum int) (*model.Withdrawal, error)
	GetByUserID(ctx context.Context, userID int) ([]model.Withdrawal, error)
}
