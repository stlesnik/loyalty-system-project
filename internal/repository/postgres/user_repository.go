package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/loyalty-system-project/internal/model"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
)

const (
	ErrCodeUniqueViolation = "23505"
)

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{db: db}
}

func (u *User) Create(ctx context.Context, login, password string) (*model.User, error) {
	var id int
	err := u.db.QueryRowContext(
		ctx,
		"INSERT INTO users(login, password) VALUES($1, $2) RETURNING id",
		login, password,
	).Scan(&id)

	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == ErrCodeUniqueViolation {
		utils.Log.Infow("User already exists", "login", login)
		return nil, utils.ErrLoginAlreadyExists
	}
	if err != nil {
		utils.Log.Infow("User creation failed", "login", login, "error", err)
		return nil, err
	}

	utils.Log.Infow("User created", "login", login, "id", id)
	return &model.User{ID: id, Login: login, PasswordHash: password}, nil
}

func (u *User) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE login = $1", login)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Log.Infow("User not found", "login", login)
		return nil, utils.ErrLoginDoesntExist
	}
	if err != nil {
		utils.Log.Infow("User get by login query failed", "login", login, "error", err)
		return nil, err
	}
	utils.Log.Infow("User found", "login", login)
	return &user, nil
}

func (u *User) GetByID(ctx context.Context, userID int) (*model.User, error) {
	var user *model.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", userID)
	if errors.Is(err, sql.ErrNoRows) {
		utils.Log.Infow("User not found", "id", userID)
		return nil, utils.ErrIDDoesntExist
	}
	if err != nil {
		utils.Log.Infow("User get by id query failed", "id", userID)
		return nil, err
	}
	utils.Log.Infow("User found", "id", userID)
	return user, nil
}

func (u *User) GetDB() *sqlx.DB { return u.db }
