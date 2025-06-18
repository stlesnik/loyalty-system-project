package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/loyalty-system-project/internal/model"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"time"
)

type Withdrawal struct {
	db *sqlx.DB
}

func NewWithdrawal(db *sqlx.DB) *Withdrawal {
	return &Withdrawal{db: db}
}

func (w *Withdrawal) Create(ctx context.Context, userID int, orderNumber string, amount float64) (*model.Withdrawal, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", 3*60*60)
	}
	processedAt := time.Now().In(loc)

	_, err = w.db.ExecContext(ctx, `
        INSERT INTO withdrawals (user_id, order_number, amount, processed_at)
        VALUES ($1, $2, $3, $4)
    `, userID, orderNumber, amount, processedAt)

	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == ErrCodeUniqueViolation {
		return nil, utils.ErrOrderAlreadyUploaded
	}
	if err != nil {
		utils.Log.Errorf("Error creating withdrawal: %s", err.Error())
		return nil, err
	}
	return &model.Withdrawal{
		UserID:      userID,
		OrderID:     orderNumber,
		Amount:      amount,
		ProcessedAt: processedAt,
	}, nil
}

func (w *Withdrawal) GetTotal(ctx context.Context, userID int) (float64, error) {
	var withdrawn float64
	err := w.db.GetContext(ctx, &withdrawn, `
        SELECT COALESCE(SUM(amount), 0) 
        FROM withdrawals 
        WHERE user_id = $1
    `, userID)
	return withdrawn, err
}

func (w *Withdrawal) GetByUserID(ctx context.Context, userID int) ([]model.Withdrawal, error) {
	var withdrawals []model.Withdrawal
	err := w.db.SelectContext(ctx, &withdrawals, `
        SELECT *
        FROM withdrawals
        WHERE user_id = $1
        ORDER BY processed_at DESC 
    `, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawals: %w", err)
	}
	return withdrawals, nil
}
