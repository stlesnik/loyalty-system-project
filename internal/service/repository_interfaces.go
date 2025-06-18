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
	Create(ctx context.Context, userID int, orderNumber string) (*model.Order, error)
	GetByUserID(ctx context.Context, userID int) ([]model.Order, error)
	GetByID(ctx context.Context, orderNumber string) (*model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

type BalanceRepository interface {
	GetTotal(ctx context.Context, userID int) (float64, error)
	Update(ctx context.Context, userID int, deltaCurrent int, deltaWithdrawn int) error
}

type WithdrawalRepository interface {
	Create(ctx context.Context, userID int, orderNumber string, sum float64) (*model.Withdrawal, error)
	GetTotal(ctx context.Context, userID int) (float64, error)
	GetByUserID(ctx context.Context, userID int) ([]model.Withdrawal, error)
}
