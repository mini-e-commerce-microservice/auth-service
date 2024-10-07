package jwt_util

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/jwt_claims_proto"
)

type AuthAccessTokenClaims struct {
	Claims *jwt_claims_proto.JwtAuthAccessTokenClaims `json:"claims"`
	jwt.RegisteredClaims
}
