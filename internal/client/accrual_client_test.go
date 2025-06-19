package client

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccrualClient_FetchAccrual(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Запускаем тестовый сервер
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"order": "123", "status": "PROCESSED", "accrual": 500}`))
			require.NoError(t, err)
		}))
		defer server.Close()

		client := NewAccrualClient(server.URL)
		resp, err := client.FetchAccrual(context.Background(), "123")

		require.NoError(t, err)
		require.Equal(t, "PROCESSED", resp.Status)
		require.Equal(t, 500.0, *resp.Accrual)
	})

	t.Run("429 retry", func(t *testing.T) {
		attempt := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if attempt == 0 {
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				attempt++
			} else {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"order": "123", "status": "PROCESSED"}`))
				require.NoError(t, err)
			}
		}))
		defer server.Close()

		client := NewAccrualClient(server.URL)
		resp, err := client.FetchAccrual(context.Background(), "123")

		require.NoError(t, err)
		require.Equal(t, "PROCESSED", resp.Status)
	})
}
