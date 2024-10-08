package auth

import (
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/token"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
)

type service struct {
	tokenRepository token.Repository
	userRepository  users.Repository
}

func NewService(tokenRepo token.Repository, userRepo users.Repository) *service {
	return &service{
		tokenRepository: tokenRepo,
		userRepository:  userRepo,
	}
}
