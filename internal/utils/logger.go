package utils

import (
	"go.uber.org/zap"
	"strings"
)

var Log *zap.SugaredLogger

func InitLogger(env string) error {
	var logger *zap.Logger
	var err error

	switch strings.ToLower(env) {
	case "prod", "production":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return err
	}
	Log = logger.Sugar()
	return nil
}
