package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
)

type Repositories struct {
	User     *User
	Order    *Order
	Balance  *Balance
	Withdraw *Withdrawal
}

func newDataBase(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		utils.Log.Errorf("error while opening db: %w: %v", utils.ErrOpenDB, err)
		return nil, fmt.Errorf("error while opening db: %w: %v", utils.ErrOpenDB, err)
	}
	return db, nil
}
func InitRepositories(dsn string) (Repositories, error) {
	db, err := newDataBase(dsn)
	if err != nil {
		return Repositories{}, err
	}

	return Repositories{
		NewUser(db),
		NewOrder(db),
		NewBalance(db),
		NewWithdrawal(db),
	}, nil
}
