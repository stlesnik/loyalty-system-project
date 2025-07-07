package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type User struct {
	ID           int    `json:"id" db:"id"`
	Login        string `json:"login" db:"login"`
	PasswordHash string `json:"password" db:"password"`
}
type Status string

const (
	StatusNew        Status = "NEW"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
)

func (s Status) Value() (driver.Value, error) {
	if s == "" {
		return nil, nil
	}
	return string(s), nil
}

func (s *Status) Scan(src interface{}) error {
	if src == nil {
		*s = StatusNew
		return nil
	}
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into Status", src)
	}
	switch Status(str) {
	case StatusNew, StatusProcessing, StatusInvalid, StatusProcessed:
		*s = Status(str)
		return nil
	}
	return fmt.Errorf("unknown status %q", str)
}

type Order struct {
	Number     string    `json:"number" db:"number"`
	UserID     int       `json:"-" db:"user_id"`
	Status     Status    `json:"status" db:"status"`
	Accrual    *float64  `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}
type Balance struct {
	UserID  int
	Current float64
}

type Withdrawal struct {
	UserID      int       `json:"-" db:"user_id"`
	OrderID     string    `json:"order" db:"order_number"`
	Amount      float64   `json:"sum" db:"amount"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

type TotalWithdrawal struct {
	UserID      int
	TotalAmount float64
}
