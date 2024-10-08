package token

import (
	"context"
	"fmt"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/tracer"
)

func (r *repository) DeleteToken(ctx context.Context, input DeleteTokenInput) (err error) {
	key := fmt.Sprintf("%s%s-%s", primitive.PrefixCacheRedis, input.TokenType, input.TokenUID)

	cmd := r.client.B().Del().Key(key).Build()

	resp := r.client.Do(ctx, cmd)
	if err = resp.Error(); err != nil {
		return tracer.Error(err)
	}
	return
}

type DeleteTokenInput struct {
	TokenType primitive.EnumTokenType
	TokenUID  string
}
