package jwt_util

import (
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/jwt_claims_proto"
	"time"
)

type AuthAccessTokenClaims struct {
	*jwt_claims_proto.JwtAuthAccessTokenClaims
	jwt.RegisteredClaims
}

func (a *AuthAccessTokenClaims) GenerateHS256(key string, expiredAt int64) (tokenStr string, err error) {
	timeNow := time.Now().UTC()
	timeExp := timeNow.Add(time.Minute * time.Duration(expiredAt))
	a.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(timeExp),
		IssuedAt:  jwt.NewNumericDate(timeNow),
		NotBefore: jwt.NewNumericDate(timeNow),
		Audience:  []string{"users"},
	}

	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, a)

	tokenStr, err = tokenParse.SignedString([]byte(key))
	if err != nil {
		return tokenStr, collection.Err(err)
	}
	return
}

func (a *AuthAccessTokenClaims) ClaimsHS256(tokenStr string, key string) (err error) {
	tokenParse, err := jwt.ParseWithClaims(tokenStr, a, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return collection.Err(err)
	}

	if !tokenParse.Valid {
		return collection.Err(ErrInvalidParseToken)
	}

	return
}
