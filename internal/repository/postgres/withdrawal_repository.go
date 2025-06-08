package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/loyalty-system-project/internal/model"
)

type Withdrawal struct {
	db *sqlx.DB
}

func NewWithdrawal(db *sqlx.DB) *Withdrawal {
	return &Withdrawal{db: db}
}
func (w *Withdrawal) CreateWithdrawal(ctx context.Context, userID int, orderID string, sum int) error {
	return nil
}
func (w *Withdrawal) GetWithdrawalsByUserID(ctx context.Context, userID int) ([]model.Withdrawal, error) {
	return nil, nil
}
