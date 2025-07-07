package main

import (
	"context"
	"errors"
	"github.com/stlesnik/loyalty-system-project/internal/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	a, err := app.New()
	if err != nil {
		panic(err)
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Сервер запущен на %s", a.Cfg.ServerAddress)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		log.Fatalf("Сервер неожиданно завершил работу: %v", err)
	case sig := <-stop:
		log.Printf("Получен сигнал %s, завершаем работу...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.Stop(ctx); err != nil {
			log.Fatalf("Ошибка при остановке приложения: %v", err)
		}
		log.Println("Сервер успешно остановлен")
	}

}
