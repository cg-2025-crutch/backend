package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Server   ServerConfig
	Logger   LoggerConfig
	Postgres PostgresConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	GRPCHost string `env:"GRPC_HOST" envDefault:""`
	GRPCPort string `env:"GRPC_PORT" envDefault:"50051"`
}

type PostgresConfig struct {
	Host           string        `env:"DB_HOST"`
	User           string        `env:"DB_USER"`
	Password       string        `env:"DB_PASSWORD"`
	Name           string        `env:"DB_NAME"`
	Port           string        `env:"DB_PORT"`
	SSLMode        string        `env:"DB_SSLMODE"`
	ConnTimeout    time.Duration `env:"DB_CONN_TIMEOUT"`
	MigrateTimeout time.Duration `env:"DB_MIGRATE_TIMEOUT"`
	MigrationsDir  string        `env:"DB_MIGRATIONS_DIR"`
}

type LoggerConfig struct {
	Development bool `env:"LOG_DEV" envDefault:"false"`
	Caller      bool `env:"LOG_CALLER" envDefault:"false"`
}

type JWTConfig struct {
	SecretKey       string        `env:"JWT_SECRET_KEY" envDefault:"your-secret-key-change-in-production"`
	AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"15m"`
	RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"168h"`
}

func New() (AppConfig, error) {
	var cfg AppConfig
	_ = godotenv.Load(".env")
	err := env.Parse(&cfg)
	if err != nil {
		return AppConfig{}, fmt.Errorf("error parsing config: %w", err)
	}
	return cfg, nil
}
