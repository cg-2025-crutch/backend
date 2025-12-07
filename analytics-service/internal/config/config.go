package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Server           ServerConfig
	Logger           LoggerConfig
	Kafka            KafkaConfig
	Redis            RedisConfig
	UserServiceAddr  string `env:"USER_SERVICE_ADDR"`
	FundsServiceAddr string `env:"FUNDS_SERVICE_ADDR"`
}

type ServerConfig struct {
	GRPCHost string `env:"GRPC_HOST" envDefault:""`
	GRPCPort string `env:"GRPC_PORT" envDefault:"50054"`
}

type LoggerConfig struct {
	Development bool `env:"LOG_DEV" envDefault:"false"`
	Caller      bool `env:"LOG_CALLER" envDefault:"false"`
}

type KafkaConfig struct {
	Brokers      []string      `env:"KAFKA_BROKERS" envDefault:"localhost:9092" envSeparator:","`
	ConsTopic    []string      `env:"KAFKA_CONS_TOPIC" envDefault:"analytics" envSeparator:","`
	ConnDeadline time.Duration `env:"KAFKA_CONN_DEADLINE" envDefault:"10s"`
}

type RedisConfig struct {
	Host         string        `env:"CACHE_HOST" envDefault:"localhost:6379"`
	Password     string        `env:"CACHE_PASSWORD" envDefault:""`
	ConnDeadline time.Duration `env:"CACHE_CONN_DEADLINE" envDefault:"10s"`
	TTL          time.Duration `env:"REDIS_TTL" envDefault:"24h"`
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
