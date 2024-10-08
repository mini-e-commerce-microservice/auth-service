package auth

import "context"

type Service interface {
	Login(ctx context.Context, input LoginInput) (output LoginOutput, err error)
	Logout(ctx context.Context, input LogoutInput) (err error)
	GenerateAccessToken(ctx context.Context, input GenerateAccessTokenInput) (output GenerateAccessTokenOutput, err error)
}
