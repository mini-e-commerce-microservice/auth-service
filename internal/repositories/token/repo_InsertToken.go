package token

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/mini-e-commerce-microservice/auth-service/internal/model"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"time"
)

func (r *repository) InsertToken(ctx context.Context, input InsertTokenInput) (err error) {
	key := fmt.Sprintf("%s%s-%s", r.redisConf.TrackingPrefix, input.TokenType, input.TokenUID)
	dataValue, err := json.Marshal(input.Value)
	if err != nil {
		return collection.Err(err)
	}

	cmd := r.client.B().Set().Key(key).Value(string(dataValue)).Exat(input.ExpiredAt).Build()

	resp := r.client.Do(ctx, cmd)
	if err = resp.Error(); err != nil {
		return collection.Err(err)
	}

	return
}

type InsertTokenInput struct {
	TokenType primitive.EnumTokenType
	TokenUID  string
	ExpiredAt time.Time
	Value     model.TokenCache
}
