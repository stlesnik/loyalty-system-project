package handler

//
//import (
//	"bytes"
//	"encoding/json"
//	"errors"
//	"github.com/stlesnik/loyalty-system-project/internal/utils"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestRegisterHandler(t *testing.T) {
//	tests := []struct {
//		name         string
//		payload      map[string]string
//		expectedCode int
//		mockError    error
//	}{
//		{
//			name:         "200 OK",
//			payload:      map[string]string{"login": "new_user", "password": "pass123"},
//			expectedCode: http.StatusOK,
//		},
//		{
//			name:         "400 Bad Request",
//			payload:      map[string]string{},
//			expectedCode: http.StatusBadRequest,
//		},
//		{
//			name:         "409 Conflict",
//			payload:      map[string]string{"login": "existing_user", "password": "pass123"},
//			expectedCode: http.StatusConflict,
//			mockError:    utils.ErrLoginExists,
//		},
//		{
//			name:         "500 Internal Error",
//			payload:      map[string]string{"login": "user", "password": "pass"},
//			expectedCode: http.StatusInternalServerError,
//			mockError:    errors.New("storage fucked up"),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			storageMock := new(MockStorage)
//			storageMock.On("CreateUser", mock.Anything, mock.Anything).Return(tt.mockError)
//
//			app := &App{storage: storageMock}
//			body, _ := json.Marshal(tt.payload)
//			req := httptest.NewRequest("POST", "/api/user/register", bytes.NewReader(body))
//			rr := httptest.NewRecorder()
//
//			app.RegisterHandler(rr, req)
//
//			assert.Equal(t, tt.expectedCode, rr.Code)
//			if tt.expectedCode == http.StatusOK {
//				assert.NotEmpty(t, rr.Result().Cookies(), "Session cookie not set")
//			}
//		})
//	}
//}
