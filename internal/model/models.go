package model

import "time"

type User struct {
	ID           int    `json:"id" db:"id"`
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"password" db:"password"`
}
type Status int

const (
	NEW Status = iota
	PROCESSING
	INVALID
	PROCESSED
)

type Order struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Status      Status    `json:"status"`
	BonusAmount int       `json:"bonus_amount"`
	CreatedAt   time.Time `json:"created_at"`
}
type Balance struct {
	UserID  int    `json:"user_id"`
	Balance string `json:"balance"`
}

type Withdrawal struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	OrderID string `json:"order_id"`
	Amount  int    `json:"amount"`
}
