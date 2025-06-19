package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stlesnik/loyalty-system-project/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthSvc struct {
	rep UserRepository
}

func NewAuthService(rep UserRepository) *AuthSvc {
	return &AuthSvc{rep: rep}
}

func (s *AuthSvc) RegisterUser(ctx context.Context, login, password string) (int, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.Log.Infow("Failed to hash password", "error", err)
		return -1, err
	}
	user, err := s.rep.Create(ctx, login, string(hashedPass))
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

func (s *AuthSvc) Authenticate(ctx context.Context, login string, password string) (int, error) {
	user, err := s.rep.GetByLogin(ctx, login)
	switch {
	case errors.Is(err, utils.ErrLoginDoesntExist):
		return -1, utils.ErrInvalidCredentials
	case err != nil:
		return -1, err
	default:
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			utils.Log.Infow("Failed to compare password", "error", err)
			return -1, utils.ErrInvalidCredentials
		}
		utils.Log.Infow("User credentials validated", "id", user.ID, "login", user.Login)
		return user.ID, nil
	}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func (s *AuthSvc) GenerateUserToken(id int, secretKey string, tokenExp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		utils.Log.Infow("Failed to sign token", "error", err)
		return "", fmt.Errorf("jwt signing failed: %w", err)
	}
	utils.Log.Infow("Token generated", "token", tokenString)
	return tokenString, nil
}
