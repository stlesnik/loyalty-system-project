package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/middleware"
	"github.com/stlesnik/loyalty-system-project/internal/model"
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"github.com/stlesnik/loyalty-system-project/internal/service/mocks"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBalanceHandler_GetBalance(t *testing.T) {
	err := utils.InitLogger("dev")
	require.NoError(t, err)
	tests := []struct {
		name           string
		setupMock      func(*mocks.MockBalanceService)
		expectedStatus int
		expectedJSON   string
	}{
		{
			name: "success",
			setupMock: func(m *mocks.MockBalanceService) {
				m.EXPECT().GetBalance(gomock.Any(), 123).
					Return(&service.Balance{
						Current:   500.5,
						Withdrawn: 42,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   `{"current":500.5,"withdrawn":42}`,
		},
		{
			name: "no balance",
			setupMock: func(m *mocks.MockBalanceService) {
				m.EXPECT().GetBalance(gomock.Any(), 123).
					Return(&service.Balance{
						Current:   0,
						Withdrawn: 0,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   `{"current":0,"withdrawn":0}`,
		},
		{
			name: "internal error",
			setupMock: func(m *mocks.MockBalanceService) {
				m.EXPECT().GetBalance(gomock.Any(), 123).
					Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockBalanceService(ctrl)
			tt.setupMock(mockSvc)

			cfg := &config.Config{}
			handler := NewBalanceHandler(mockSvc, cfg)

			req := httptest.NewRequest("GET", "/api/user/balance", nil)
			ctx := context.WithValue(req.Context(), middleware.UserIDKeyName, 123)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.GetBalance(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedJSON != "" {
				require.JSONEq(t, tt.expectedJSON, rr.Body.String())
			}
		})
	}
}

func TestBalanceHandler_Withdraw(t *testing.T) {
	err := utils.InitLogger("dev")
	require.NoError(t, err)

	tests := []struct {
		name           string
		body           map[string]interface{}
		setupMock      func(*mocks.MockBalanceService)
		expectedStatus int
	}{
		{
			name: "success",
			body: map[string]interface{}{
				"order": "2377225624",
				"sum":   751.0,
			},
			setupMock: func(m *mocks.MockBalanceService) {
				m.EXPECT().CreateWithdrawal(gomock.Any(), 123, "2377225624", 751.0).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid json",
			body:           nil,
			setupMock:      func(m *mocks.MockBalanceService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid order format",
			body: map[string]interface{}{
				"order": "123",
				"sum":   100.0,
			},
			setupMock:      func(m *mocks.MockBalanceService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "insufficient funds",
			body: map[string]interface{}{
				"order": "2377225624",
				"sum":   1000.0,
			},
			setupMock: func(m *mocks.MockBalanceService) {
				m.EXPECT().CreateWithdrawal(gomock.Any(), 123, "2377225624", 1000.0).
					Return(utils.ErrInsufficientFunds)
			},
			expectedStatus: http.StatusPaymentRequired,
		},
		{
			name: "order conflict",
			body: map[string]interface{}{
				"order": "2377225624",
				"sum":   500.0,
			},
			setupMock: func(m *mocks.MockBalanceService) {
				m.EXPECT().CreateWithdrawal(gomock.Any(), 123, "2377225624", 500.0).
					Return(utils.ErrOrderAlreadyUploaded)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "internal error",
			body: map[string]interface{}{
				"order": "2377225624",
				"sum":   300.0,
			},
			setupMock: func(m *mocks.MockBalanceService) {
				m.EXPECT().CreateWithdrawal(gomock.Any(), 123, "2377225624", 300.0).
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockBalanceService(ctrl)
			tt.setupMock(mockSvc)

			cfg := &config.Config{}
			handler := NewBalanceHandler(mockSvc, cfg)

			var body []byte
			if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			} else {
				body = []byte("{invalid-json}")
			}

			req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			ctx := context.WithValue(req.Context(), middleware.UserIDKeyName, 123)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.Withdraw(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestBalanceHandler_GetWithdrawals(t *testing.T) {
	err := utils.InitLogger("dev")
	require.NoError(t, err)

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockBalanceService)
		expectedStatus int
		expectedJSON   string
	}{
		{
			name: "success with data",
			setupMock: func(m *mocks.MockBalanceService) {
				withdrawals := []model.Withdrawal{
					{UserID: 123, OrderID: "order1", Amount: 500.0, ProcessedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
					{UserID: 123, OrderID: "order2", Amount: 300.0, ProcessedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
				}
				m.EXPECT().GetAllWithdrawals(gomock.Any(), 123).Return(withdrawals, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON: `[
                {"order": "order1", "sum": 500, "processed_at": "2023-01-01T00:00:00Z"},
                {"order": "order2", "sum": 300, "processed_at": "2023-01-02T00:00:00Z"}
            ]`,
		},
		{
			name: "no data",
			setupMock: func(m *mocks.MockBalanceService) {
				m.EXPECT().GetAllWithdrawals(gomock.Any(), 123).Return(nil, nil)
			},
			expectedStatus: http.StatusNoContent, // 204
		},
		{
			name: "internal error",
			setupMock: func(m *mocks.MockBalanceService) {
				m.EXPECT().GetAllWithdrawals(gomock.Any(), 123).Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockBalanceService(ctrl)
			tt.setupMock(mockSvc)

			cfg := &config.Config{}
			handler := NewBalanceHandler(mockSvc, cfg)

			req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
			ctx := context.WithValue(req.Context(), middleware.UserIDKeyName, 123)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.GetWithdrawals(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedJSON != "" {
				require.JSONEq(t, tt.expectedJSON, rr.Body.String())
			}
		})
	}
}
