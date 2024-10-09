package infra

import (
	"github.com/mini-e-commerce-microservice/auth-service/generated/proto/secret_proto"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

func NewMinio(cred *secret_proto.Minio) *minio.Client {
	minioClient, err := minio.New(cred.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cred.AccessId, cred.SecretAccessKey, ""),
		Secure: cred.UseSsl,
	})
	if err != nil {
		panic(err)
	}

	log.Info().Msg("initialization minio successfully")
	return minioClient
}
