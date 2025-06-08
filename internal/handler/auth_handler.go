package handler

import (
	"github.com/stlesnik/loyalty-system-project/internal/service"
	"net/http"
)

type AuthHandler struct {
	s *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{s: s} }

func (a *AuthHandler) Register(res http.ResponseWriter, req *http.Request) {}
func (a *AuthHandler) Login(res http.ResponseWriter, req *http.Request)    {}
