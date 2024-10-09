package auth

import (
	"context"
	"errors"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/token"
	jwt_util "github.com/mini-e-commerce-microservice/auth-service/internal/util/jwt"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/tracer"
)

// Logout
// list return error: ErrInvalidToken
func (s *service) Logout(ctx context.Context, input LogoutInput) (err error) {
	refreshTokenClaims := jwt_util.AuthRefreshTokenClaims{}
	err = refreshTokenClaims.ClaimsHS256(input.RefreshToken, s.jwtConf.RefreshToken.Key)
	if err != nil {
		err = errors.Join(err, ErrInvalidToken)
		return tracer.Error(err)
	}

	err = s.tokenRepository.DeleteToken(ctx, token.DeleteTokenInput{
		TokenType: primitive.EnumTokenTypeRT,
		TokenUID:  refreshTokenClaims.Uid,
	})
	if err != nil {
		return tracer.Error(err)
	}
	return
}

type LogoutInput struct {
	RefreshToken string
}
