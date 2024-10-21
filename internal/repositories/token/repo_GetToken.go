package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/mini-e-commerce-microservice/auth-service/internal/model"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"github.com/redis/rueidis"
	"time"
)

func (r *repository) GetToken(ctx context.Context, input GetTokenInput) (output GetTokenOutput, err error) {
	if input.TimeToLiveCache == 0 {
		input.TimeToLiveCache = 24 * time.Hour
	}

	key := fmt.Sprintf("%s%s-%s", r.redisConf.TrackingPrefix, input.TokenType, input.TokenUID)

	cmd := r.client.B().Get().Key(key).Cache()
	cache := r.client.DoCache(ctx, cmd, input.TimeToLiveCache)

	data, err := cache.ToString()
	if err != nil {
		if errors.Is(err, rueidis.Nil) {
			err = repositories.ErrRecordNotFound
		}
		return output, collection.Err(err)
	}

	if err = json.Unmarshal([]byte(data), &output.Data); err != nil {
		return output, collection.Err(err)
	}

	return
}

type GetTokenInput struct {
	TokenType       primitive.EnumTokenType
	TokenUID        string
	TimeToLiveCache time.Duration
}

type GetTokenOutput struct {
	Data model.TokenCache
}
