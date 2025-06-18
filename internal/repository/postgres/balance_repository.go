package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Balance struct {
	db *sqlx.DB
}

func NewBalance(db *sqlx.DB) *Balance {
	return &Balance{db: db}
}

func (b *Balance) GetTotal(ctx context.Context, userID int) (float64, error) {
	var current float64
	err := b.db.GetContext(ctx, &current, `
        SELECT COALESCE(SUM(accrual), 0) 
        FROM orders 
        WHERE user_id = $1 AND status = 'PROCESSED'
    `, userID)
	return current, err
}

func (b *Balance) Update(ctx context.Context, userID int, deltaCurrent int, deltaWithdrawn int) error {
	return nil
}
