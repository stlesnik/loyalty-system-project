package model

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

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

var statusStrings = [...]string{
	"NEW",
	"PROCESSING",
	"INVALID",
	"PROCESSED",
}

func (s *Status) String() string {
	if s == nil {
		return "<nil>"
	}
	if int(*s) < 0 || int(*s) >= len(statusStrings) {
		return "UNKNOWN"
	}
	return statusStrings[*s]
}

func (s *Status) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return s.String(), nil
}

func (s *Status) Scan(value interface{}) error {
	if value == nil {
		*s = NEW
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("status: cannot scan non-string %T", value)
	}
	switch strVal {
	case "NEW":
		*s = NEW
	case "PROCESSING":
		*s = PROCESSING
	case "INVALID":
		*s = INVALID
	case "PROCESSED":
		*s = PROCESSED
	default:
		return fmt.Errorf("status: unknown value %q", strVal)
	}
	return nil
}

type Order struct {
	Number     string        `json:"number" db:"number"`
	UserID     int           `json:"-" db:"user_id"`
	Status     Status        `json:"status" db:"status"`
	Accrual    sql.NullInt64 `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time     `json:"uploaded_at" db:"uploaded_at"`
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
