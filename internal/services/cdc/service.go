package cdc

import (
	ekafka "github.com/SyaibanAhmadRamadhan/event-bus/kafka"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/secret_proto"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type cdc struct {
	kafkaBroker    ekafka.KafkaPubSub
	userRepository users.Repository
	kafkaConf      *secret_proto.Kafka
	propagators    propagation.TextMapPropagator
	dbTransaction  wsqlx.Tx
}

func New(kafkaBroker ekafka.KafkaPubSub, kafkaConf *secret_proto.Kafka, userRepository users.Repository, dbTransaction wsqlx.Tx) *cdc {
	return &cdc{
		userRepository: userRepository,
		propagators:    otel.GetTextMapPropagator(),
		kafkaBroker:    kafkaBroker,
		dbTransaction:  dbTransaction,
		kafkaConf:      kafkaConf,
	}
}
