package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/tracer"
	"github.com/redis/rueidis"
	"time"
)

func (r *repository) CheckToken(ctx context.Context, input CheckTokenInput) (exists bool, err error) {
	if input.TimeToLiveCache == 0 {
		input.TimeToLiveCache = 24 * time.Hour
	}

	key := fmt.Sprintf("%s%s-%s", primitive.PrefixCacheRedis, input.TokenType, input.TokenUID)

	cmd := r.client.B().Get().Key(key).Cache()
	cache := r.client.DoCache(ctx, cmd, input.TimeToLiveCache)

	_, err = cache.ToString()
	if err != nil {
		if errors.Is(err, rueidis.Nil) {
			return false, nil
		}
		return exists, tracer.Error(err)
	}

	return true, nil
}

type CheckTokenInput struct {
	TokenType       primitive.EnumTokenType
	TokenUID        string
	TimeToLiveCache time.Duration
}
