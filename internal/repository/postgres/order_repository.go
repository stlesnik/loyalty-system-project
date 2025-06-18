package postgres

import (
	"context"
	"database/sql"
	"errors"
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
	var id int
	createdAt := time.Now().UTC()

	err := o.db.QueryRowContext(
		ctx,
		"INSERT INTO orders(user_id, order_id, created_at) VALUES($1, $2, $3) RETURNING id",
		userID, orderID, createdAt,
	).Scan(&id)

	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == ErrCodeUniqueViolation {
		utils.Log.Infow("Order already exists", "order_id", orderID)
		return nil, utils.ErrLoginAlreadyExists
	}
	if err != nil {
		utils.Log.Infow("Order creation failed", "order_id", orderID, "error", err)
		return nil, err
	}

	utils.Log.Infow("Order created", "order_id", orderID, "id", id)
	return &model.Order{
		ID:          id,
		UserID:      userID,
		Status:      model.NEW,
		BonusAmount: 0,
		CreatedAt:   createdAt,
	}, nil
}

func (o *Order) GetByUserID(ctx context.Context, userID int) ([]model.Order, error) {
	return nil, nil
}

func (o *Order) GetByID(ctx context.Context, orderID string) (*model.Order, error) {
	var order *model.Order
	err := o.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE order_id = $1", orderID)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Log.Infow("Order not found", "id", orderID)
		return nil, utils.ErrIDDoesntExist
	}
	if err != nil {
		utils.Log.Infow("Order get by id query failed", "id", orderID)
		return nil, err
	}
	utils.Log.Infow("Order found", "id", orderID)
	return order, nil
}

func (o *Order) Update(ctx context.Context, order model.Order) error { return nil }
