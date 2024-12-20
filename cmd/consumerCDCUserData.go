package main

import (
	"context"
	ekafka "github.com/SyaibanAhmadRamadhan/event-bus/kafka"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/auth-service/internal/conf"
	"github.com/mini-e-commerce-microservice/auth-service/internal/infra"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
	"github.com/mini-e-commerce-microservice/auth-service/internal/services/cdc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
)

var consumerCDCUserData = &cobra.Command{
	Use:   "consumerCDCUserData",
	Short: "consumerCDCUserData",
	Run: func(cmd *cobra.Command, args []string) {
		kafkaConf := conf.LoadKafkaConf()
		appConf := conf.LoadAppConf()
		otelConf := conf.LoadOtelConf()

		db, closeFnPostgre := infra.NewPostgresql(appConf.DatabaseDsn)
		closeFnOtel := infra.NewOtel(otelConf, appConf.TracerName)
		rdbms := wsqlx.NewRdbms(db)
		kafkaBroker := ekafka.New(ekafka.WithOtel())

		userRepository := users.NewRepository(rdbms)

		cdcService := cdc.New(kafkaBroker, kafkaConf, userRepository, rdbms)

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			err := cdcService.ConsumerUserData(ctx)
			if err != nil {
				log.Err(err)
				stop()
			}
		}()

		<-ctx.Done()
		log.Info().Msg("Received shutdown signal, shutting down server gracefully...")

		closeFnOtel(context.TODO())
		closeFnPostgre(context.TODO())

		log.Info().Msg("Shutdown complete. Exiting.")
	},
}
