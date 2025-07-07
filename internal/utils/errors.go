package utils

import (
	"errors"
)

var (
	ErrOpenDB        = errors.New("error while opening db")
	ErrIDDoesntExist = errors.New("id does not exist")

	//auth
	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrLoginDoesntExist   = errors.New("login does not exist")
	ErrInvalidCredentials = errors.New("user not found")

	//order
	ErrOrderAlreadyUploaded = errors.New("order already uploaded")
	ErrOrderConflict        = errors.New("order uploaded by another user")
	ErrNoOrders             = errors.New("no orders found")

	//balance
	ErrInsufficientFunds = errors.New("insufficient funds")

	//accruel
	ErrOrderNotFoundInAccrual = errors.New("order not found")
)
