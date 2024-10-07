package cdc

import "context"

type Service interface {
	ConsumerUserData(ctx context.Context) (err error)
}
