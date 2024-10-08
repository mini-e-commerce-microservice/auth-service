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
		JwtAuthAccessTokenClaims: &jwt_claims_proto.JwtAuthAccessTokenClaims{
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
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0NTIxNjE3ODkyNTcyMTQ3Njg0LCJlbWFpbCI6InBJWkl3ZmFASmdvY0htQi5uZXQiLCJyZWdpc3Rlcl9hcyI6MSwiaXNzIjoidGVzdCIsInN1YiI6InNvbWVib2R5IiwiYXVkIjpbInNvbWVib2R5X2Vsc2UiXSwiZXhwIjoxNzI4NDA0NjU3LCJuYmYiOjE3MjgzMTgyNTcsImlhdCI6MTcyODMxODI1NywianRpIjoiMSJ9.03QvgCaK_DAsbRslIZxilQJrdhggUyxptabFjpGDskk"

	type claimsToken AuthAccessTokenClaims

	token, err := jwt.ParseWithClaims(tokenStr, &claimsToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		panic(err)
	} else if claims, ok := token.Claims.(*claimsToken); ok {
		fmt.Println(claims.UserId, claims.RegisteredClaims.Issuer)
	} else {
		panic("unknown claims type, cannot proceed")
	}
}
