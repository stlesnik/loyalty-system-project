package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stlesnik/loyalty-system-project/internal/config"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"net/http"
)

type contextKey string

const (
	UserIDKeyName contextKey = "userID"
)

func RequireAuth(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromCookie(r, cfg.AuthSecretKey)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		utils.Log.Infow("Got user id from cookie", "userID", userID)
		ctx := context.WithValue(r.Context(), UserIDKeyName, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func getUserIDFromCookie(r *http.Request, secretKey string) (int, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return -1, fmt.Errorf("failed to get Authorization cookie")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return -1, err
	}
	if !token.Valid {
		return -1, errors.New("invalid token")
	}

	return claims.UserID, nil
}
