package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/jwt_claims_proto"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/token"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
	jwt_util "github.com/mini-e-commerce-microservice/auth-service/internal/util/jwt"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/tracer"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Login
// list of error type: ErrInvalidEmail, ErrInvalidPassword
func (s *service) Login(ctx context.Context, input LoginInput) (output LoginOutput, err error) {
	userOutput, err := s.userRepository.FindOneUser(ctx, users.FindOneUserInput{
		Email: null.StringFrom(input.Email),
	})
	if err != nil {
		if !errors.Is(err, repositories.ErrRecordNotFound) {
			return output, tracer.Error(err)
		}
		return output, tracer.Error(ErrInvalidEmail)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userOutput.Data.Password), []byte(input.Password))
	if err != nil {
		return output, tracer.Error(ErrInvalidPassword)
	}

	accessTokenClaim := jwt_util.AuthAccessTokenClaims{
		JwtAuthAccessTokenClaims: &jwt_claims_proto.JwtAuthAccessTokenClaims{
			UserId:     userOutput.Data.ID,
			Email:      userOutput.Data.Email,
			RegisterAs: int64(userOutput.Data.RegisterAs),
		},
	}
	accessToken, err := accessTokenClaim.GenerateHS256(s.jwtConf.AccessToken.Key, s.jwtConf.AccessToken.ExpiredAt)
	if err != nil {
		return output, tracer.Error(err)
	}

	refreshTokenClaim := jwt_util.AuthRefreshTokenClaims{
		JwtAuthRefreshTokenClaims: &jwt_claims_proto.JwtAuthRefreshTokenClaims{
			UserId:     userOutput.Data.ID,
			Email:      userOutput.Data.Email,
			RegisterAs: int64(userOutput.Data.RegisterAs),
			Uid:        uuid.New().String(),
		},
	}
	expiredAt := s.jwtConf.RefreshToken.ExpiredAt
	if input.RememberMe {
		expiredAt = s.jwtConf.RefreshToken.RememberMeExpiredAt
	}

	refreshToken, err := refreshTokenClaim.GenerateHS256(s.jwtConf.RefreshToken.Key, expiredAt)
	if err != nil {
		return output, tracer.Error(err)
	}

	err = s.tokenRepository.InsertToken(ctx, token.InsertTokenInput{
		TokenType: primitive.EnumTokenTypeRT,
		Token:     refreshToken,
		TokenUID:  refreshTokenClaim.Uid,
		ExpiredAt: refreshTokenClaim.ExpiresAt.UTC(),
	})
	if err != nil {
		return output, tracer.Error(err)
	}

	_, err = s.tokenRepository.CheckToken(ctx, token.CheckTokenInput{
		TokenType:       primitive.EnumTokenTypeRT,
		TokenUID:        refreshTokenClaim.Uid,
		TimeToLiveCache: 48 * time.Hour,
	})
	if err != nil {
		return output, tracer.Error(err)
	}

	output = LoginOutput{
		AccessToken: LoginOutputToken{
			Token:     accessToken,
			ExpiredAt: accessTokenClaim.ExpiresAt.UTC(),
		},
		RefreshToken: LoginOutputToken{
			Token:     refreshToken,
			ExpiredAt: refreshTokenClaim.ExpiresAt.UTC(),
		},
		User: LoginOutputUser{
			ID:              userOutput.Data.ID,
			Email:           userOutput.Data.Email,
			IsEmailVerified: userOutput.Data.IsEmailVerified,
		},
	}
	return
}

type LoginInput struct {
	Email      string
	Password   string
	RememberMe bool
}

type LoginOutput struct {
	AccessToken  LoginOutputToken
	RefreshToken LoginOutputToken
	User         LoginOutputUser
}

type LoginOutputToken struct {
	Token     string
	ExpiredAt time.Time
}

type LoginOutputUser struct {
	ID              int64
	Email           string
	IsEmailVerified bool
}
