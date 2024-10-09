package infra

import (
	"context"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/secret_proto"
	"github.com/mini-e-commerce-microservice/auth-service/internal/util/primitive"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidisotel"
	"github.com/rs/zerolog/log"
)

func NewRedisWithOtel(redisConf *secret_proto.Redis, clientName string) (rueidis.Client, primitive.CloseFn) {
	client, err := rueidisotel.NewClient(rueidis.ClientOption{
		Password:              redisConf.Password,
		ClientName:            clientName,
		InitAddress:           []string{redisConf.Host},
		ClientTrackingOptions: []string{"PREFIX", redisConf.TrackingPrefix, "BCAST"},
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
