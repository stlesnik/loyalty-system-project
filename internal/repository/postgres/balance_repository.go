package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/loyalty-system-project/internal/model"
)

type Balance struct {
	db *sqlx.DB
}

func NewBalance(db *sqlx.DB) *Balance {
	return &Balance{db: db}
}
func (b *Balance) Get(ctx context.Context, userID int) (*model.Balance, error) {
	return nil, nil
}
func (b *Balance) Update(ctx context.Context, userID int, deltaCurrent int, deltaWithdrawn int) error {
	return nil
}
