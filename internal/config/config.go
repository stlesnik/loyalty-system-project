package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"time"
)

type Config struct {
	ServerAddress        string `env:"RUN_ADDRESS"`
	DatabaseDSN          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`

	Environment   string        `env:"ENVIRONMENT"`
	AuthSecretKey string        `env:"AUTH_SECRET_KEY"`
	AuthTokenExp  time.Duration `env:"AUTH_TOKEN_EXP"`
}

func New() (*Config, error) {
	cfg := &Config{}

	defaultAddress := "localhost:8080"
	defaultDatabaseDSN := ""
	defaultAccrualSystemAddress := ""
	defaultEnvironment := "dev"
	defaultAuthSecretKey := "loyalty_system_secret_key"
	defaultAuthTokenExp := time.Hour * 24

	flag.StringVar(&cfg.ServerAddress, "a", defaultAddress, "Address to run the server")
	flag.StringVar(&cfg.DatabaseDSN, "d", defaultDatabaseDSN, "Database url")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", defaultAccrualSystemAddress, "Accrual system address")
	flag.StringVar(&cfg.Environment, "e", defaultEnvironment, "Environment")
	flag.StringVar(&cfg.AuthSecretKey, "s", defaultAuthSecretKey, "Secret key for jwt token generation")
	flag.DurationVar(&cfg.AuthTokenExp, "t", defaultAuthTokenExp, "Token expiration time")
	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
