package auth

import (
	"context"
	"errors"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/jwt_claims_proto"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/token"
	jwt_util "github.com/mini-e-commerce-microservice/auth-service/internal/util/jwt"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/tracer"
	"time"
)

// GenerateAccessToken
// list return error: ErrInvalidToken, ErrRefreshTokenNotExistsInRedis
func (s *service) GenerateAccessToken(ctx context.Context, input GenerateAccessTokenInput) (output GenerateAccessTokenOutput, err error) {
	refreshTokenClaims := jwt_util.AuthRefreshTokenClaims{}
	err = refreshTokenClaims.ClaimsHS256(input.RefreshToken)
	if err != nil {
		err = errors.Join(err, ErrInvalidToken)
		return output, tracer.Error(err)
	}

	refreshTokenExists, err := s.tokenRepository.CheckToken(ctx, token.CheckTokenInput{
		TokenType:       primitive.EnumTokenTypeRT,
		TokenUID:        refreshTokenClaims.Uid,
		TimeToLiveCache: 48 * time.Hour,
	})
	if err != nil {
		return output, tracer.Error(err)
	}
	if !refreshTokenExists {
		return output, tracer.Error(ErrRefreshTokenNotExistsInRedis)
	}

	accessTokenClaim := jwt_util.AuthAccessTokenClaims{
		JwtAuthAccessTokenClaims: &jwt_claims_proto.JwtAuthAccessTokenClaims{
			UserId:     refreshTokenClaims.UserId,
			Email:      refreshTokenClaims.Email,
			RegisterAs: refreshTokenClaims.RegisterAs,
		},
	}
	accessToken, err := accessTokenClaim.GenerateHS256()
	if err != nil {
		return output, tracer.Error(err)
	}

	output = GenerateAccessTokenOutput{
		Token:     accessToken,
		ExpiredAt: accessTokenClaim.ExpiresAt.UTC(),
	}
	return
}

type GenerateAccessTokenInput struct {
	RefreshToken string
}

type GenerateAccessTokenOutput struct {
	Token     string
	ExpiredAt time.Time
}
