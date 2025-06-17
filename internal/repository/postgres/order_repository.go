package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/loyalty-system-project/internal/model"
)

type Order struct {
	db *sqlx.DB
}

func NewOrder(db *sqlx.DB) *Order {
	return &Order{db: db}
}

func (o *Order) CreateOrder(ctx context.Context, userID int, orderID string) (*model.Order, error) {
	return nil, nil
}
func (o *Order) GetOrdersByUserID(ctx context.Context, userID int) ([]model.Order, error) {
	return nil, nil
}
func (o *Order) GetOrderByID(ctx context.Context, orderID string) (*model.Order, error) {
	return nil, nil
}
func (o *Order) UpdateOrder(ctx context.Context, order model.Order) error { return nil }
