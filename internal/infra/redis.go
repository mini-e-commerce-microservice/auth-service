package infra

import (
	"context"
	"github.com/mini-e-commerce-microservice/auth-service/internal/conf"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisotel"
	"github.com/rs/zerolog/log"
)

func NewRedisWithOtel(redisConf conf.ConfigRedis) (rueidis.Client, primitive.CloseFn) {
	client, err := rueidisotel.NewClient(rueidis.ClientOption{
		Password:              redisConf.Password,
		ClientName:            redisConf.ClientName,
		InitAddress:           []string{redisConf.Host},
		ClientTrackingOptions: []string{"PREFIX", primitive.PrefixCacheRedis, "BCAST"},
		CacheSizeEachConn:     rueidis.DefaultCacheBytes,
	})
	if err != nil {
		panic(err)
	}

	fn := func(ctx context.Context) error {
		log.Info().Msg("start closed redis client...")
		client.Close()
		log.Info().Msg("closed redis client successfully")
		return nil
	}

	return client, fn
}
