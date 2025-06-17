package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/loyalty-system-project/internal/model"
)

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{db: db}
}
func (u *User) CreateUser(ctx context.Context, login string, password string) error { return nil }

func (u *User) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	return nil, nil
}
func (u *User) GetUserByID(ctx context.Context, userID int) (*model.User, error) { return nil, nil }

func (u *User) GetDB() *sqlx.DB { return u.db }
