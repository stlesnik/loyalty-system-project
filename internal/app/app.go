package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/handler"
	"net/http"
)

type App struct {
	router chi.Router
	cfg    *config.Config
}

func New(authH *handler.AuthHandler, ordH *handler.OrderHandler, balH *handler.BalanceHandler, config *config.Config) *App {
	a := &App{
		router: chi.NewRouter(),
		cfg:    config,
	}
	a.initRouter(authH, ordH, balH)
	return a
}

func (a *App) Start() error {
	return http.ListenAndServe(a.cfg.ServerAddress, a.router)
}
