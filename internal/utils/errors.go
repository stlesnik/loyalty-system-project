package utils

import "errors"

var (
	ErrOpenDB = errors.New("error while opening db")
	ErrWarmDB = errors.New("error while warming db up")
	ErrPingDB = errors.New("error while ping to db")
)
