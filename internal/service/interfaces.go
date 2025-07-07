package service

import (
	"context"
	"github.com/stlesnik/loyalty-system-project/internal/model"
	"time"
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
	Update(ctx context.Context, orderNumber string, status string, accrual *float64) error
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

type OrderService interface {
	StartWorkers(ctx context.Context, workerCount int)
	UploadOrder(ctx context.Context, userID int, orderNumber string) error
	GetUserOrders(ctx context.Context, userID int) ([]model.Order, error)
}
type BalanceService interface {
	GetBalance(ctx context.Context, userID int) (*Balance, error)
	CreateWithdrawal(ctx context.Context, userID int, orderNumber string, amount float64) error
	GetAllWithdrawals(ctx context.Context, userID int) ([]model.Withdrawal, error)
}
type AuthService interface {
	RegisterUser(ctx context.Context, login, password string) (int, error)
	Authenticate(ctx context.Context, login string, password string) (int, error)
	GenerateUserToken(id int, secretKey string, tokenExp time.Duration) (string, error)
}
