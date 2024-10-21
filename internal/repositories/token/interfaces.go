package token

import "context"

type Repository interface {
	InsertToken(ctx context.Context, input InsertTokenInput) (err error)
	GetToken(ctx context.Context, input GetTokenInput) (output GetTokenOutput, err error)
	DeleteToken(ctx context.Context, input DeleteTokenInput) (err error)
}
