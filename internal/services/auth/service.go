package auth

import (
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/secret_proto"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/token"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
)

type service struct {
	tokenRepository token.Repository
	userRepository  users.Repository
	jwtConf         *secret_proto.Jwt
}

func NewService(tokenRepo token.Repository, userRepo users.Repository, jwtConf *secret_proto.Jwt) *service {
	return &service{
		tokenRepository: tokenRepo,
		userRepository:  userRepo,
		jwtConf:         jwtConf,
	}
}
