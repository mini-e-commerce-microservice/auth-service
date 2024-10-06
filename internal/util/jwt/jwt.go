package jwt_util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/tracer"
	"time"
)

type Jwt struct {
	Key     string
	Exp     time.Duration
	Payload string
}

func GenerateHS256(jwtModel Jwt) (string, error) {
	timeNow := time.Now()
	timeExp := timeNow.Add(jwtModel.Exp).Unix()

	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     timeExp,
		"payload": jwtModel.Payload,
	})

	tokenStr, err := tokenParse.SignedString([]byte(jwtModel.Key))
	if err != nil {
		return "", tracer.Error(err)
	}

	return tokenStr, nil
}

func ClaimHS256(token, key string) (map[string]any, error) {
	tokenParse, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, tracer.Error(fmt.Errorf("unexpected signing method: %v", t.Header["alg"]))
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, tracer.Error(err)
	}

	claims, _ := tokenParse.Claims.(jwt.MapClaims)

	return claims, nil
}
