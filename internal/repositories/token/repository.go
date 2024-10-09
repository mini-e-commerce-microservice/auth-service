package token

import (
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/secret_proto"
	"github.com/redis/rueidis"
)

type repository struct {
	client    rueidis.Client
	redisConf *secret_proto.Redis
}

func NewRepository(client rueidis.Client, redisConf *secret_proto.Redis) *repository {
	return &repository{
		client:    client,
		redisConf: redisConf,
	}
}
