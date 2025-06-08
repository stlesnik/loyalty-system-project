package service

import "github.com/stlesnik/loyalty-system-project/internal/repository"

type AuthService struct {
	rep repository.UserRepository
}

func NewAuthService(rep repository.UserRepository) *AuthService {
	return &AuthService{rep: rep}
}
