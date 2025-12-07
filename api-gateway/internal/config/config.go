package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Server                 ServerConfig
	UserServiceClient      ServiceClientConfig
	FundsServiceClient     ServiceClientConfig
	NotifServiceClient     ServiceClientConfig
	AnalyticsServiceClient ServiceClientConfig
}

type ServerConfig struct {
	Host string `env:"API_GATEWAY_HOST" envDefault:""`
	Port string `env:"API_GATEWAY_PORT" envDefault:"8080"`
}

type ServiceClientConfig struct {
	Host    string        `env:"HOST"`
	Port    string        `env:"PORT"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"10s"`
}

func New() (AppConfig, error) {
	var cfg AppConfig
	_ = godotenv.Load(".env")

	// Parse main config
	err := env.Parse(&cfg.Server)
	if err != nil {
		return AppConfig{}, fmt.Errorf("error parsing server config: %w", err)
	}

	// Parse user service config
	err = env.ParseWithOptions(&cfg.UserServiceClient, env.Options{
		Prefix: "USER_SERVICE_",
	})
	if err != nil {
		return AppConfig{}, fmt.Errorf("error parsing user service config: %w", err)
	}

	// Parse funds service config
	err = env.ParseWithOptions(&cfg.FundsServiceClient, env.Options{
		Prefix: "FUNDS_SERVICE_",
	})
	if err != nil {
		return AppConfig{}, fmt.Errorf("error parsing funds service config: %w", err)
	}

	// Parse notification service config
	err = env.ParseWithOptions(&cfg.NotifServiceClient, env.Options{
		Prefix: "NOTIFICATION_SERVICE_",
	})
	if err != nil {
		return AppConfig{}, fmt.Errorf("error parsing notification service config: %w", err)
	}

	// Parse analytics service config
	err = env.ParseWithOptions(&cfg.AnalyticsServiceClient, env.Options{
		Prefix: "ANALYTICS_SERVICE_",
	})
	if err != nil {
		return AppConfig{}, fmt.Errorf("error parsing analytics service config: %w", err)
	}

	return cfg, nil
}
