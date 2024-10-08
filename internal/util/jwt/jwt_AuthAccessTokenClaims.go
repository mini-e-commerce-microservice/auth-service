package jwt_util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/jwt_claims_proto"
	"github.com/mini-e-commerce-microservice/auth-service/internal/conf"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/tracer"
	"time"
)

type AuthAccessTokenClaims struct {
	*jwt_claims_proto.JwtAuthAccessTokenClaims
	jwt.RegisteredClaims
}

func (a *AuthAccessTokenClaims) GenerateHS256() (tokenStr string, err error) {
	accessTokenConf := conf.GetConfig().Jwt.AccessToken

	timeNow := time.Now().UTC()
	timeExp := timeNow.Add(accessTokenConf.ExpiredAt)
	a.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(timeExp),
		IssuedAt:  jwt.NewNumericDate(timeNow),
		NotBefore: jwt.NewNumericDate(timeNow),
		Audience:  []string{"users"},
	}

	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, a)

	tokenStr, err = tokenParse.SignedString([]byte(accessTokenConf.Key))
	if err != nil {
		return tokenStr, tracer.Error(err)
	}
	return
}

func (a *AuthAccessTokenClaims) ClaimsHS256(tokenStr string) (err error) {
	accessTokenConf := conf.GetConfig().Jwt.AccessToken

	tokenParse, err := jwt.ParseWithClaims(tokenStr, a, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessTokenConf.Key), nil
	})
	if err != nil {
		return tracer.Error(err)
	}

	if !tokenParse.Valid {
		return tracer.Error(ErrInvalidParseToken)
	}

	return
}
