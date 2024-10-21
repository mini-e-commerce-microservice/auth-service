package token

import (
	"context"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
)

func (r *repository) DeleteToken(ctx context.Context, input DeleteTokenInput) (err error) {
	key := fmt.Sprintf("%s%s-%s", r.redisConf.TrackingPrefix, input.TokenType, input.TokenUID)

	cmd := r.client.B().Del().Key(key).Build()

	resp := r.client.Do(ctx, cmd)
	if err = resp.Error(); err != nil {
		return collection.Err(err)
	}
	return
}

type DeleteTokenInput struct {
	TokenType primitive.EnumTokenType
	TokenUID  string
}
