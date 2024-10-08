package main

import (
	"context"
	"errors"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/auth-service/internal/conf"
	"github.com/mini-e-commerce-microservice/auth-service/internal/infra"
	"github.com/mini-e-commerce-microservice/auth-service/internal/presenter"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/token"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
	"github.com/mini-e-commerce-microservice/auth-service/internal/services/auth"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"net/http"
	"os/signal"
	"syscall"
)

var restApi = &cobra.Command{
	Use:   "rest-api",
	Short: "REST API",
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init()
		config := conf.GetConfig()

		redisClient, redisClose := infra.NewRedisWithOtel(config.Redis)
		otelClose := infra.NewOtel(config.OpenTelemetry)
		postgreClient, postgreClose := infra.NewPostgresql(config.DatabaseDSN)
		rdbms := wsqlx.NewRdbms(postgreClient, wsqlx.WithAttributes(
			semconv.DBSystemPostgreSQL,
		))

		tokenRepository := token.NewRepository(redisClient)
		usersRepository := users.NewRepository(rdbms)

		authService := auth.NewService(tokenRepository, usersRepository)

		server := presenter.New(&presenter.Presenter{
			AuthService: authService,
			Port:        config.AppPort,
		})

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		go func() {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				ctx.Done()
			}
		}()

		<-ctx.Done()
		log.Info().Msg("Received shutdown signal, shutting down server gracefully...")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Err(err).Msg("failed shutdown server")
		}

		//time.Sleep(30 * time.Second)

		redisClose(context.Background())

		if err := otelClose(context.Background()); err != nil {
			log.Err(err).Msg("failed closed otel")
		}

		if err := postgreClose(context.Background()); err != nil {
			log.Err(err).Msg("failed closed postgresql")
		}

		log.Info().Msg("Shutdown complete. Exiting.")
	},
}
