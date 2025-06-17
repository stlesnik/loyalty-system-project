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
	a.router.Post("/api/user/register", wrap(authH.Register))        // регистрация пользователя;
	a.router.Post("/api/user/login", wrap(authH.Login))              // аутентификация пользователя;
	a.router.Post("/api/user/orders", wrap(ordH.UploadOrder))        // загрузка пользователем номера заказа для расчёта;
	a.router.Get("/api/user/orders", wrap(ordH.GetOrders))           // получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	a.router.Get("/api/user/balance", wrap(balH.GetBalance))         // получение текущего баланса счёта баллов лояльности пользователя;
	a.router.Post("/api/user/balance/withdraw", wrap(balH.Withdraw)) // запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	a.router.Get("/api/user/withdrawals", wrap(balH.GetWithdrawals)) // получение информации о выводе средств с накопительного счёта пользователем.
}
