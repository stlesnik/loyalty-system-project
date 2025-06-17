package service

type AuthService struct {
	rep UserRepository
}

func NewAuthService(rep UserRepository) *AuthService {
	return &AuthService{rep: rep}
}
