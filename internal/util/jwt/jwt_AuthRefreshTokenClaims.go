package jwt_util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/jwt_claims_proto"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/tracer"
	"time"
)

type AuthRefreshTokenClaims struct {
	*jwt_claims_proto.JwtAuthRefreshTokenClaims
	jwt.RegisteredClaims
}

func (a *AuthRefreshTokenClaims) GenerateHS256(key string, expiredAt int64) (tokenStr string, err error) {
	timeNow := time.Now().UTC()
	timeExp := timeNow.Add(time.Duration(expiredAt) * time.Minute)

	a.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(timeExp),
		IssuedAt:  jwt.NewNumericDate(timeNow),
		NotBefore: jwt.NewNumericDate(timeNow),
		Audience:  []string{"users"},
	}

	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, a)

	tokenStr, err = tokenParse.SignedString([]byte(key))
	if err != nil {
		return tokenStr, tracer.Error(err)
	}
	return
}

func (a *AuthRefreshTokenClaims) ClaimsHS256(tokenStr string, key string) (err error) {
	tokenParse, err := jwt.ParseWithClaims(tokenStr, a, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return tracer.Error(err)
	}

	if !tokenParse.Valid {
		return tracer.Error(ErrInvalidParseToken)
	}

	return
}
