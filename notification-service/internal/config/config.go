package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Notif  NotificationsConfig
	Redis  RedisConfig
	Server ServerConfig
	Logger LoggerConfig
	Kafka  KafkaConfig
}

type RedisConfig struct {
	Host         string        `env:"CACHE_HOST" envDefault:"cache:6379"`
	Password     string        `env:"CACHE_PASSWORD"`
	ConnDeadline time.Duration `env:"CACHE_CONN_DEADLINE" envDefault:"10s"`
}

type ServerConfig struct {
	GRPCHost string `env:"GRPC_HOST" envDefault:""`
	GRPCPort string `env:"GRPC_PORT" envDefault:"50051"`
	HTTPHost string `env:"HTTP_HOST" envDefault:""`
	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`
}

type LoggerConfig struct {
	Development bool `env:"LOG_DEV" envDefault:"false"`
	Caller      bool `env:"LOG_CALLER" envDefault:"false"`
}

type KafkaConfig struct {
	Brokers      []string      `env:"KAFKA_BROKERS" env-required:"true"`
	ConsTopic    []string      `env:"KAFKA_CONS_TOPIC" env-required:"true"`
	ProdTopic    string        `env:"KAFKA_PROD_TOPIC" env-required:"true"`
	ConnDeadline time.Duration `env:"KAFKA_CONN_DEADLINE" envDefault:"10s"`
}

type NotificationsConfig struct {
	VapidPublic  string `env:"VAPID_PUBLIC" envDefault:""`
	VapidPrivate string `env:"VAPID_PRIVATE" envDefault:""`
	Subscriber   string `env:"VAPID_SUBSCRIBER" envDefault:"mailto:admin@example.com"`
}

func New() (AppConfig, error) {
	var cfg AppConfig
	_ = godotenv.Load()
	err := env.Parse(&cfg)
	if err != nil {
		return AppConfig{}, fmt.Errorf("error parsing config: %w", err)
	}
	return cfg, nil
}
