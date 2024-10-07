package conf

import "time"

type Config struct {
	AppPort       int                 `mapstructure:"APP_PORT"`
	AppMode       string              `mapstructure:"APP_MODE"`
	OpenTelemetry ConfigOpenTelemetry `mapstructure:"OPEN_TELEMETRY"`
	Minio         ConfigMinio         `mapstructure:"MINIO"`
	DatabaseDSN   string              `mapstructure:"DATABASE_DSN"`
	Kafka         ConfigKafka         `mapstructure:"KAFKA"`
	Jwt           ConfigJWT           `mapstructure:"JWT"`
}

type ConfigKafka struct {
	Host  string           `mapstructure:"HOST"`
	Topic ConfigKafkaTopic `mapstructure:"TOPIC"`
}

type ConfigKafkaTopic struct {
	CDCUserTable string `mapstructure:"TOPIC_CDC_USER_TABLE"`
}

type ConfigOpenTelemetry struct {
	Password   string `mapstructure:"PASSWORD"`
	Username   string `mapstructure:"USERNAME"`
	Endpoint   string `mapstructure:"ENDPOINT"`
	TracerName string `mapstructure:"TRACER_NAME"`
}

type ConfigMinio struct {
	Endpoint        string `mapstructure:"ENDPOINT"`
	AccessID        string `mapstructure:"ACCESS_ID"`
	SecretAccessKey string `mapstructure:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `mapstructure:"USE_SSL"`
	PrivateBucket   string `mapstructure:"PRIVATE_BUCKET"`
}

type ConfigJWT struct {
	AccessToken  ConfigJWTAccessToken  `mapstructure:"ACCESS_TOKEN"`
	RefreshToken ConfigJWTRefreshToken `mapstructure:"REFRESH_TOKEN"`
}

type ConfigJWTAccessToken struct {
	Key       string        `mapstructure:"KEY"`
	ExpiredAt time.Duration `mapstructure:"EXPIRED_AT"`
}

type ConfigJWTRefreshToken struct {
	Key                 string        `mapstructure:"KEY"`
	ExpiredAt           time.Duration `mapstructure:"EXPIRED_AT"`
	RememberMeExpiredAt time.Duration `mapstructure:"REMEMBER_ME_EXPIRED_AT"`
}
