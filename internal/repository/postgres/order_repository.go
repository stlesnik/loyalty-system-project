package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/loyalty-system-project/internal/model"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"time"
)

type Order struct {
	db *sqlx.DB
}

func NewOrder(db *sqlx.DB) *Order {
	return &Order{db: db}
}

func (o *Order) Create(ctx context.Context, userID int, orderID string) (*model.Order, error) {
	uploadedAt := time.Now().UTC()

	_, err := o.db.ExecContext(
		ctx,
		"INSERT INTO orders(user_id, number, uploaded_at) VALUES($1, $2, $3)",
		userID, orderID, uploadedAt,
	)

	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == ErrCodeUniqueViolation {
		utils.Log.Infow("Order already exists", "number", orderID)
		return nil, utils.ErrLoginAlreadyExists
	}
	if err != nil {
		utils.Log.Infow("Order creation failed", "number", orderID, "error", err)
		return nil, err
	}

	utils.Log.Infow("Order uploaded", "number", orderID)
	return &model.Order{
		Number:     orderID,
		UserID:     userID,
		Status:     model.NEW,
		UploadedAt: uploadedAt,
	}, nil
}

func (o *Order) GetByUserID(ctx context.Context, userID int) ([]model.Order, error) {
	var orders []model.Order
	err := o.db.SelectContext(ctx, &orders, `
        SELECT *
        FROM orders 
        WHERE user_id = $1
        ORDER BY created_at DESC  
    `, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (o *Order) GetByID(ctx context.Context, orderID string) (*model.Order, error) {
	var order model.Order
	err := o.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE number = $1", orderID)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Log.Infow("Order not found", "id", orderID)
		return nil, utils.ErrIDDoesntExist
	}
	if err != nil {
		utils.Log.Errorw("Order get by id query failed", "error", err, "id", orderID)
		return nil, err
	}
	utils.Log.Infow("Order found", "id", orderID)
	return &order, nil
}

func (o *Order) Update(ctx context.Context, order model.Order) error { return nil }
