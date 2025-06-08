package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress        string `env:"RUN_ADDRESS"`
	DatabaseDSN          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`

	Environment   string `env:"ENVIRONMENT"`
	AuthSecretKey string `env:"AUTH_SECRET_KEY"`
}

func New() (*Config, error) {
	cfg := &Config{}

	defaultAddress := "localhost:8080"
	defaultDatabaseDSN := ""
	defaultAccrualSystemAddress := ""
	defaultEnvironment := "dev"
	defaultAuthSecretKey := "url_shortener_secret_key"

	flag.StringVar(&cfg.ServerAddress, "a", defaultAddress, "Address to run the server")
	flag.StringVar(&cfg.DatabaseDSN, "d", defaultDatabaseDSN, "Database url")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", defaultAccrualSystemAddress, "Accrual system address")
	flag.StringVar(&cfg.Environment, "e", defaultEnvironment, "Environment")
	flag.StringVar(&cfg.AuthSecretKey, "s", defaultAuthSecretKey, "Secret key for jwt token generation")
	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
