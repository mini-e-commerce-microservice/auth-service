package cdc

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SyaibanAhmadRamadhan/event-bus/debezium"
	ekafka "github.com/SyaibanAhmadRamadhan/event-bus/kafka"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/auth-service/internal/model"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"go.opentelemetry.io/otel/trace"
)

func (c *cdc) ConsumerUserData(ctx context.Context) (err error) {
	output, err := c.kafkaBroker.Subscribe(ctx, ekafka.SubInput{
		Config: kafka.ReaderConfig{
			Brokers: []string{c.kafkaConf.Host},
			GroupID: "user-service-consumer-user-data-group1",
			Topic:   c.kafkaConf.Topic.CdcUserTable,
		},
	})
	if err != nil {
		return collection.Err(err)
	}

	for {
		data := DebeziumPayload[UserData]{}
		msg, err := output.Reader.FetchMessage(ctx, &data)
		if err != nil {
			if !errors.Is(err, ekafka.ErrJsonUnmarshal) {
				return collection.Err(err)
			}
			continue
		}

		carrier := ekafka.NewMsgCarrier(&msg)
		ctxConsumer := c.propagators.Extract(context.Background(), carrier)

		ctxConsumer, span := otel.Tracer("").Start(ctxConsumer, string(data.Payload.Op)+" process cdc user data from user service.",
			trace.WithAttributes(
				attribute.String("cdc.debezium.payload.op", string(data.Payload.Op)),
				attribute.Int64("cdc.debezium.payload.data.id", data.Payload.ID),
				attribute.String("cdc.debezium.payload.data.email", data.Payload.Email),
				attribute.Int64("cdc.debezium.payload.data.register_as", int64(data.Payload.RegisterAs)),
				attribute.Bool("cdc.debezium.payload.data.is_email_verified", data.Payload.IsEmailVerified),
			))

		switch data.Payload.Op {
		case debezium.Create, debezium.Update:
			_ = c.dbTransaction.DoTxContext(ctxConsumer, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  false,
			}, func(ctx context.Context, tx wsqlx.Rdbms) error {
				_, err = c.userRepository.UpSertUser(ctx, users.UpSertUserInput{
					Tx: tx,
					Payload: model.User{
						ID:              data.Payload.ID,
						Email:           data.Payload.Email,
						Password:        data.Payload.Password,
						CreatedAt:       data.Payload.CreatedAt,
						IsEmailVerified: data.Payload.IsEmailVerified,
						RegisterAs:      data.Payload.RegisterAs,
						UpdatedAt:       data.Payload.UpdatedAt,
						DeletedAt:       data.Payload.DeletedAt,
						TraceParent:     data.Payload.TraceParent,
					},
				})
				if err != nil {
					span.RecordError(collection.Err(err))
					span.SetStatus(codes.Error, err.Error())
					span.SetAttributes(semconv.ErrorTypeKey.String("failed create user"))
					return err
				}

				err = output.Reader.CommitMessages(ctx, msg)
				if err != nil {
					span.RecordError(collection.Err(err))
					span.SetStatus(codes.Error, err.Error())
					span.SetAttributes(semconv.ErrorTypeKey.String("failed commit message"))
					return err
				}

				span.SetStatus(codes.Ok, "cdc successfully")
				return nil
			})
		default:
			err = output.Reader.CommitMessages(ctx, msg)
			if err != nil {
				span.RecordError(collection.Err(err))
				span.SetStatus(codes.Error, err.Error())
				span.SetAttributes(semconv.ErrorTypeKey.String("failed commit message"))
				return err
			}
			span.SetStatus(codes.Error, "unsupported debezium operation type")
			span.SetAttributes(semconv.ErrorTypeKey.String("unsupported debezium operation type"))
		}
		span.End()
	}

}
