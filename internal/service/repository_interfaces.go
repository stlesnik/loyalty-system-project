package service

import (
	"context"
	"github.com/stlesnik/loyalty-system-project/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login, password string) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	GetUserByID(ctx context.Context, userID int) (*model.User, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, userID int, orderID string) (*model.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int) ([]model.Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order model.Order) error
}

type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int) (*model.Balance, error)
	UpdateBalance(ctx context.Context, userID int, deltaCurrent int, deltaWithdrawn int) error
}

type WithdrawalRepository interface {
	CreateWithdrawal(ctx context.Context, userID int, orderID string, sum int) (*model.Withdrawal, error)
	GetWithdrawalsByUserID(ctx context.Context, userID int) ([]model.Withdrawal, error)
}
