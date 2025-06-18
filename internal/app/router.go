package app

import (
	"github.com/stlesnik/loyalty-system-project/internal/handler"
	"github.com/stlesnik/loyalty-system-project/internal/middleware"
	"net/http"
)

func (a *App) initRouter(authH *handler.AuthHandler, ordH *handler.OrderHandler, balH *handler.BalanceHandler) {
	wrap := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.WithLogging(
			middleware.WithDecompress(
				middleware.WithCompress(h),
			),
		)
	}
	authWrap := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.RequireAuth(a.Cfg, wrap(h))
	}

	a.router.Post("/api/user/register", wrap(authH.Register))
	a.router.Post("/api/user/login", wrap(authH.Login))
	a.router.Post("/api/user/orders", authWrap(ordH.UploadOrder))
	a.router.Get("/api/user/orders", wrap(ordH.GetOrders))
	a.router.Get("/api/user/balance", wrap(balH.GetBalance))
	a.router.Post("/api/user/balance/withdraw", wrap(balH.Withdraw))
	a.router.Get("/api/user/withdrawals", wrap(balH.GetWithdrawals))
}
