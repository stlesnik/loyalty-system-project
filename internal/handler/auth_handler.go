package handler

import (
	"encoding/json"
	"errors"
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"net/http"
	"time"
)

type AuthHandler struct {
	s   *service.AuthService
	cfg *config.Config
}

func NewAuthHandler(s *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{s: s, cfg: cfg}
}

func (aH *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID, err := aH.s.RegisterUser(r.Context(), req.Login, req.Password)
	switch {
	case errors.Is(err, utils.ErrLoginAlreadyExists):
		http.Error(w, "Login taken", http.StatusConflict)
	case err != nil:
		http.Error(w, "Server error", http.StatusInternalServerError)
	default:
		tokenString, _ := aH.s.GenerateUserToken(userID, aH.cfg.AuthSecretKey, aH.cfg.AuthTokenExp)

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Expires:  time.Now().Add(aH.cfg.AuthTokenExp),
			HttpOnly: true,
			Path:     "/",
		})
		w.WriteHeader(http.StatusOK)
	}
}
func (aH *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID, err := aH.s.Authenticate(r.Context(), req.Login, req.Password)
	switch {
	case errors.Is(err, utils.ErrInvalidCredentials):
		http.Error(w, "Invalid login/password", http.StatusUnauthorized)
	case err != nil:
		http.Error(w, "Server error", http.StatusInternalServerError)
	default:
		tokenString, _ := aH.s.GenerateUserToken(userID, aH.cfg.AuthSecretKey, aH.cfg.AuthTokenExp)

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Expires:  time.Now().Add(aH.cfg.AuthTokenExp),
			HttpOnly: true,
			Path:     "/",
		})
		w.WriteHeader(http.StatusOK)
	}
}
