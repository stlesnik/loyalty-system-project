package main

import (
	"fmt"
	"github.com/stlesnik/loyalty-system-project/internal/app"
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/handler"
	"github.com/stlesnik/loyalty-system-project/internal/repository/postgres"
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"log"
)

func main() {
	// config
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Не получилось обработать конфиг: %s", err)
		return
	}

	//logger
	err = utils.InitLogger(cfg.Environment)
	if err != nil {
		panic(fmt.Errorf("logger broke: %w", err))
	}

	//repository
	reps, err := postgres.InitRepositories(cfg.DatabaseDSN)
	if err != nil {
		panic(fmt.Errorf("db could not start: %w", err))
	}

	//services
	authSvc := service.NewAuthService(reps.User)
	orderSvc := service.NewOrderService(reps.Order)
	balanceSvc := service.NewBalanceService(reps.Balance, reps.Withdraw)

	// handlers
	authHandler := handler.NewAuthHandler(authSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	balanceHandler := handler.NewBalanceHandler(balanceSvc)

	// app
	a := app.New(
		authHandler,
		orderHandler,
		balanceHandler,
		cfg,
	)

	log.Printf("Сервер запущен на %s", cfg.ServerAddress)
	err = a.Start()
	if err != nil {
		log.Fatalf("Не получилось запустить приложение: %s", err)
		return
	}
}
