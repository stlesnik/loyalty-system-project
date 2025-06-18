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

func (o *Order) Create(ctx context.Context, userID int, orderNumber string) (*model.Order, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", 3*60*60)
	}
	uploadedAt := time.Now().In(loc)

	_, err = o.db.ExecContext(
		ctx,
		"INSERT INTO orders(user_id, number, uploaded_at) VALUES($1, $2, $3)",
		userID, orderNumber, uploadedAt,
	)

	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == ErrCodeUniqueViolation {
		utils.Log.Infow("Order already exists", "number", orderNumber)
		return nil, utils.ErrLoginAlreadyExists
	}
	if err != nil {
		utils.Log.Infow("Order creation failed", "number", orderNumber, "error", err)
		return nil, err
	}

	utils.Log.Infow("Order uploaded", "number", orderNumber)
	return &model.Order{
		Number:     orderNumber,
		UserID:     userID,
		Status:     model.StatusNew,
		UploadedAt: uploadedAt,
	}, nil
}

func (o *Order) GetByUserID(ctx context.Context, userID int) ([]model.Order, error) {
	var orders []model.Order
	err := o.db.SelectContext(ctx, &orders, `
        SELECT *
        FROM orders 
        WHERE user_id = $1
        ORDER BY "uploaded_at" DESC  
    `, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (o *Order) GetByID(ctx context.Context, orderNumber string) (*model.Order, error) {
	var order model.Order
	err := o.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE number = $1", orderNumber)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Log.Infow("Order not found", "id", orderNumber)
		return nil, utils.ErrIDDoesntExist
	}
	if err != nil {
		utils.Log.Errorw("Order get by id query failed", "error", err, "id", orderNumber)
		return nil, err
	}
	utils.Log.Infow("Order found", "id", orderNumber)
	return &order, nil
}

func (o *Order) Update(ctx context.Context, orderNumber string, status string, accrual *float64) error {
	var dbAccrual sql.NullFloat64

	if accrual != nil {
		dbAccrual = sql.NullFloat64{
			Float64: *accrual,
			Valid:   true,
		}
	} else {
		dbAccrual = sql.NullFloat64{Valid: false}
	}

	var updatedNumber string
	err := o.db.GetContext(ctx, &updatedNumber, `
        UPDATE orders 
        SET 
            status = $1,
            accrual = $2
        WHERE number = $3
        RETURNING number
    `, status, dbAccrual, orderNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrNoOrders
		}
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}
