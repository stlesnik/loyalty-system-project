package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/handler"
	"github.com/stlesnik/loyalty-system-project/internal/repository/postgres"
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"log"
	"net/http"
)

type App struct {
	router chi.Router
	Cfg    *config.Config
	server *http.Server
	reps   postgres.Repositories
}

func New() (*App, error) {
	// config
	cfg, err := config.New()
	if err != nil {
		log.Printf("Не получилось обработать конфиг: %s", err)
		return nil, err
	}

	//logger
	err = utils.InitLogger(cfg.Environment)
	if err != nil {
		return nil, fmt.Errorf("logger broke: %w", err)
	}

	//repository
	reps, err := postgres.InitRepositories(cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("db could not start: %w", err)
	}

	//services
	authSvc := service.NewAuthService(reps.User)
	orderSvc := service.NewOrderService(reps.Order)
	balanceSvc := service.NewBalanceService(reps.Balance, reps.Withdraw)

	// handlers
	authHandler := handler.NewAuthHandler(authSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	balanceHandler := handler.NewBalanceHandler(balanceSvc)

	a := &App{
		router: chi.NewRouter(),
		Cfg:    cfg,
		server: &http.Server{},
		reps:   reps,
	}
	a.initRouter(authHandler, orderHandler, balanceHandler)
	return a, nil
}

func (a *App) Start() error {
	a.server.Handler = a.router
	a.server.Addr = a.Cfg.ServerAddress
	return a.server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.reps.User.GetDB().Close(); err != nil {
		return err
	}

	return nil
}
