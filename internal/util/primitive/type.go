package primitive

import "context"

type CloseFn func(ctx context.Context) (err error)
