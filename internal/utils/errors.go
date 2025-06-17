package utils

import (
	"errors"
)

var (
	ErrOpenDB = errors.New("error while opening db")
	ErrWarmDB = errors.New("error while warming db up")
	ErrPingDB = errors.New("error while ping to db")

	ErrLoginAlreadyExists = errors.New("login already exists")
	ErrLoginDoesntExist   = errors.New("login does not exist")
	ErrIDDoesntExist      = errors.New("id does not exist")
	ErrInvalidCredentials = errors.New("user not found")
)
