package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/middleware"
	"github.com/stlesnik/loyalty-system-project/internal/model"
	"github.com/stlesnik/loyalty-system-project/internal/service/mocks"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOrderHandler_UploadOrder(t *testing.T) {
	err := utils.InitLogger("dev")
	require.NoError(t, err)
	tests := []struct {
		name           string
		body           string
		contentType    string
		setupMock      func(*mocks.MockOrderService)
		expectedStatus int
	}{
		{
			name:        "success",
			body:        "12345678903",
			contentType: "text/plain",
			setupMock: func(m *mocks.MockOrderService) {
				m.EXPECT().UploadOrder(gomock.Any(), 123, "12345678903").Return(nil)
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "invalid luhn",
			body:           "123",
			contentType:    "text/plain",
			setupMock:      func(m *mocks.MockOrderService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:        "already uploaded",
			body:        "12345678903",
			contentType: "text/plain",
			setupMock: func(m *mocks.MockOrderService) {
				m.EXPECT().UploadOrder(gomock.Any(), 123, "12345678903").
					Return(utils.ErrOrderAlreadyUploaded)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "conflict",
			body:        "12345678903",
			contentType: "text/plain",
			setupMock: func(m *mocks.MockOrderService) {
				m.EXPECT().UploadOrder(gomock.Any(), 123, "12345678903").
					Return(utils.ErrOrderConflict)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "internal error",
			body:        "12345678903",
			contentType: "text/plain",
			setupMock: func(m *mocks.MockOrderService) {
				m.EXPECT().UploadOrder(gomock.Any(), 123, "12345678903").
					Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockOrderService(ctrl)
			tt.setupMock(mockSvc)
			cfg := &config.Config{}
			handler := NewOrderHandler(mockSvc, cfg)

			req := httptest.NewRequest("POST", "/api/user/orders", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			ctx := context.WithValue(req.Context(), middleware.UserIDKeyName, 123)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.UploadOrder(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestOrderHandler_GetOrders(t *testing.T) {
	err := utils.InitLogger("dev")
	require.NoError(t, err)

	now := time.Now().UTC().Truncate(time.Second)
	sum := 500.0
	testOrders := []model.Order{
		{
			Number:     "123",
			Status:     "PROCESSED",
			Accrual:    &sum,
			UploadedAt: now,
		},
		{
			Number:     "456",
			Status:     "PROCESSING",
			UploadedAt: now.Add(-time.Hour),
		},
	}

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockOrderService)
		expectedStatus int
		expectedJSON   string
	}{
		{
			name: "success with orders",
			setupMock: func(m *mocks.MockOrderService) {
				m.EXPECT().GetUserOrders(gomock.Any(), 123).
					Return(testOrders, nil)
			},
			expectedStatus: http.StatusOK,
			expectedJSON: ToJSON(t, []map[string]interface{}{
				{
					"number":      "123",
					"status":      "PROCESSED",
					"accrual":     500.0,
					"uploaded_at": now,
				},
				{
					"number":      "456",
					"status":      "PROCESSING",
					"uploaded_at": now.Add(-time.Hour).Truncate(time.Second),
				},
			}),
		},
		{
			name: "no orders",
			setupMock: func(m *mocks.MockOrderService) {
				m.EXPECT().GetUserOrders(gomock.Any(), 123).
					Return(nil, utils.ErrNoOrders)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "internal error",
			setupMock: func(m *mocks.MockOrderService) {
				m.EXPECT().GetUserOrders(gomock.Any(), 123).
					Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockOrderService(ctrl)
			tt.setupMock(mockSvc)

			cfg := &config.Config{}
			handler := NewOrderHandler(mockSvc, cfg)

			req := httptest.NewRequest("GET", "/api/user/orders", nil)
			ctx := context.WithValue(req.Context(), middleware.UserIDKeyName, 123)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.GetOrders(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedJSON != "" {
				require.JSONEq(t, tt.expectedJSON, rr.Body.String())
			}
		})
	}
}

func ToJSON(t *testing.T, v interface{}) string {
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return string(data)
}
