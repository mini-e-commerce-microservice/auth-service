package jwt_util

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/jwt_claims_proto"
	"math/rand"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	mySigningKey := []byte("AllYourBase")

	claims := AuthAccessTokenClaims{
		Claims: &jwt_claims_proto.JwtAuthAccessTokenClaims{
			UserId:     rand.Int63(),
			Email:      faker.Email(),
			RegisterAs: 1,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Println(ss, err)
}

func TestValidateToken(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbXMiOnsidXNlcl9pZCI6ODQ4MjYwMDkzMDM5NDc1NTkyNywiZW1haWwiOiJaanN2bm12QFpBdXVCaXguYml6IiwicmVnaXN0ZXJfYXMiOjF9LCJpc3MiOiJ0ZXN0Iiwic3ViIjoic29tZWJvZHkiLCJhdWQiOlsic29tZWJvZHlfZWxzZSJdLCJleHAiOjE3Mjg0MDQyMzksIm5iZiI6MTcyODMxNzgzOSwiaWF0IjoxNzI4MzE3ODM5LCJqdGkiOiIxIn0.GP44njB54jFL2dUN3vbcqaZ0B1m3ihosOaePvaRKBWE"

	type claimsToken AuthAccessTokenClaims

	token, err := jwt.ParseWithClaims(tokenStr, &claimsToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		panic(err)
	} else if claims, ok := token.Claims.(*claimsToken); ok {
		fmt.Println(claims.Claims, claims.RegisteredClaims.Issuer)
	} else {
		panic("unknown claims type, cannot proceed")
	}
}
