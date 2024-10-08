package token

import "context"

type Repository interface {
	InsertToken(ctx context.Context, input InsertTokenInput) (err error)
	CheckToken(ctx context.Context, input CheckTokenInput) (exists bool, err error)
	DeleteToken(ctx context.Context, input DeleteTokenInput) (err error)
}
