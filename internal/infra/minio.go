package infra

import (
	"github.com/mini-e-commerce-microservice/auth-service/internal/conf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

func NewMinio(cred conf.ConfigMinio) *minio.Client {
	minioClient, err := minio.New(cred.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cred.AccessID, cred.SecretAccessKey, ""),
		Secure: cred.UseSSL,
	})
	if err != nil {
		panic(err)
	}

	log.Info().Msg("initialization minio successfully")
	return minioClient
}
