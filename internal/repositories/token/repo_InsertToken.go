package token

import (
	"context"
	"fmt"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/tracer"
	"time"
)

func (r *repository) InsertToken(ctx context.Context, input InsertTokenInput) (err error) {
	key := fmt.Sprintf("%s%s-%s", primitive.PrefixCacheRedis, input.TokenType, input.TokenUID)

	cmd := r.client.B().Set().Key(key).Value(input.Token).Exat(input.ExpiredAt).Build()

	resp := r.client.Do(ctx, cmd)
	if err = resp.Error(); err != nil {
		return tracer.Error(err)
	}

	return
}

type InsertTokenInput struct {
	TokenType primitive.EnumTokenType
	Token     string
	TokenUID  string
	ExpiredAt time.Time
}
