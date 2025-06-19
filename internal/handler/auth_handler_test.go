package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/service/mocks"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Register(t *testing.T) {
	err := utils.InitLogger("dev")
	require.NoError(t, err)

	tests := []struct {
		name           string
		body           map[string]string
		setupMock      func(*mocks.MockAuthService)
		expectedStatus int
		expectCookie   bool
	}{
		{
			name: "success",
			body: map[string]string{"login": "test", "password": "pass123"},
			setupMock: func(m *mocks.MockAuthService) {
				m.EXPECT().RegisterUser(gomock.Any(), "test", "pass123").
					Return(123, nil)
				m.EXPECT().GenerateUserToken(123, "secret", 24*time.Hour).
					Return("valid_token", nil)
			},
			expectedStatus: http.StatusOK,
			expectCookie:   true,
		},
		{
			name: "login taken",
			body: map[string]string{"login": "taken", "password": "pass"},
			setupMock: func(m *mocks.MockAuthService) {
				m.EXPECT().RegisterUser(gomock.Any(), "taken", "pass").
					Return(0, utils.ErrLoginAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "invalid json",
			body:           nil, // Специально невалидный JSON
			setupMock:      func(m *mocks.MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			body: map[string]string{"login": "test", "password": "pass"},
			setupMock: func(m *mocks.MockAuthService) {
				m.EXPECT().RegisterUser(gomock.Any(), "test", "pass").
					Return(0, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockAuthService(ctrl)
			tt.setupMock(mockSvc)

			cfg := &config.Config{
				AuthSecretKey: "secret",
				AuthTokenExp:  24 * time.Hour,
			}

			authHandler := NewAuthHandler(mockSvc, cfg)

			var body []byte
			if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			} else {
				body = []byte("{invalid-json}")
			}

			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			authHandler.Register(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectCookie {
				cookies := rr.Result().Cookies()
				require.NotEmpty(t, cookies)

				var authCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == "auth_token" {
						authCookie = cookie
						break
					}
				}

				require.NotNil(t, authCookie)
				require.Equal(t, "valid_token", authCookie.Value)
				require.True(t, authCookie.Expires.After(time.Now()))
				require.True(t, authCookie.HttpOnly)
				require.Equal(t, "/", authCookie.Path)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	err := utils.InitLogger("dev")
	require.NoError(t, err)

	tests := []struct {
		name           string
		body           map[string]string
		setupMock      func(*mocks.MockAuthService)
		expectedStatus int
		expectCookie   bool
	}{
		{
			name: "success",
			body: map[string]string{"login": "valid", "password": "correct"},
			setupMock: func(m *mocks.MockAuthService) {
				m.EXPECT().Authenticate(gomock.Any(), "valid", "correct").
					Return(123, nil)
				m.EXPECT().GenerateUserToken(123, "secret", 24*time.Hour).
					Return("valid_token", nil)
			},
			expectedStatus: http.StatusOK,
			expectCookie:   true,
		},
		{
			name: "invalid credentials",
			body: map[string]string{"login": "user", "password": "wrong"},
			setupMock: func(m *mocks.MockAuthService) {
				m.EXPECT().Authenticate(gomock.Any(), "user", "wrong").
					Return(0, utils.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "internal error",
			body: map[string]string{"login": "test", "password": "pass"},
			setupMock: func(m *mocks.MockAuthService) {
				m.EXPECT().Authenticate(gomock.Any(), "test", "pass").
					Return(0, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSvc := mocks.NewMockAuthService(ctrl)
			tt.setupMock(mockSvc)

			cfg := &config.Config{
				AuthSecretKey: "secret",
				AuthTokenExp:  24 * time.Hour,
			}

			authHandler := NewAuthHandler(mockSvc, cfg)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			authHandler.Login(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectCookie {
				cookies := rr.Result().Cookies()
				require.NotEmpty(t, cookies)
				require.Equal(t, "auth_token", cookies[0].Name)
				require.Equal(t, "valid_token", cookies[0].Value)
			}
		})
	}
}
