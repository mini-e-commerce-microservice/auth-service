package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/jwt_claims_proto"
	"github.com/mini-e-commerce-microservice/auth-service/internal/model"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/token"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
	jwt_util "github.com/mini-e-commerce-microservice/auth-service/internal/util/jwt"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"time"
)

// GenerateAccessToken
// list return error: ErrInvalidToken, ErrRefreshTokenNotExistsInRedis
func (s *service) GenerateAccessToken(ctx context.Context, input GenerateAccessTokenInput) (output GenerateAccessTokenOutput, err error) {
	refreshTokenClaims := jwt_util.AuthRefreshTokenClaims{}
	err = refreshTokenClaims.ClaimsHS256(input.RefreshToken, s.jwtConf.RefreshToken.Key)
	if err != nil {
		err = errors.Join(err, ErrInvalidToken)
		return output, collection.Err(err)
	}

	if input.UpdateUserDataInCache {
		userOutput, err := s.userRepository.FindOneUser(ctx, users.FindOneUserInput{
			ID: null.IntFrom(refreshTokenClaims.UserId),
		})
		if err != nil {
			if errors.Is(err, repositories.ErrRecordNotFound) {
				err = fmt.Errorf("user not found in database, %w", ErrInvalidToken)
			}
			return output, collection.Err(err)
		}

		err = s.tokenRepository.InsertToken(ctx, token.InsertTokenInput{
			TokenType: primitive.EnumTokenTypeRT,
			TokenUID:  refreshTokenClaims.Uid,
			ExpiredAt: refreshTokenClaims.ExpiresAt.UTC(),
			Value: model.TokenCache{
				Email:           userOutput.Data.Email,
				IsEmailVerified: userOutput.Data.IsEmailVerified,
				RegisterAs:      userOutput.Data.RegisterAs,
			},
		})
		if err != nil {
			return output, collection.Err(err)
		}
	}

	refreshTokenData, err := s.tokenRepository.GetToken(ctx, token.GetTokenInput{
		TokenType:       primitive.EnumTokenTypeRT,
		TokenUID:        refreshTokenClaims.Uid,
		TimeToLiveCache: 48 * time.Hour,
	})
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			err = ErrRefreshTokenNotExistsInRedis
		}
		return output, collection.Err(err)
	}

	accessTokenClaim := jwt_util.AuthAccessTokenClaims{
		JwtAuthAccessTokenClaims: &jwt_claims_proto.JwtAuthAccessTokenClaims{
			UserId:          refreshTokenClaims.UserId,
			Email:           refreshTokenData.Data.Email,
			RegisterAs:      int64(refreshTokenData.Data.RegisterAs),
			IsEmailVerified: refreshTokenData.Data.IsEmailVerified,
		},
	}
	accessToken, err := accessTokenClaim.GenerateHS256(s.jwtConf.AccessToken.Key, s.jwtConf.AccessToken.ExpiredAt)
	if err != nil {
		return output, collection.Err(err)
	}

	output = GenerateAccessTokenOutput{
		Token:     accessToken,
		ExpiredAt: accessTokenClaim.ExpiresAt.UTC(),
	}
	return
}

type GenerateAccessTokenInput struct {
	RefreshToken          string
	UpdateUserDataInCache bool
}

type GenerateAccessTokenOutput struct {
	Token     string
	ExpiredAt time.Time
}
