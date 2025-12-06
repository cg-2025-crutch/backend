package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/cg-2025-crutch/backend/notification-service/internal/config"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/log"
	"github.com/redis/rueidis"
	"go.uber.org/zap"
)

func NewRedisClient(ctx context.Context, conf config.RedisConfig) (*rueidis.Client, error) {
	l := log.FromContext(ctx)

	cl, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{conf.Host},
		Password:    conf.Password,
	})
	if err == nil {
		l.Infof("connected to redis on host: %v ", conf.Host)

		return &cl, nil
	}
	l.Error("failed to connect to redis, trying 5 seconds to reconnect", zap.Error(err))
	ticker := time.NewTicker(1 * time.Second)
	ctxStop, cancel := context.WithTimeout(ctx, conf.ConnDeadline)
	defer cancel()

	for {
		select {
		case <-ctxStop.Done():
			return nil, fmt.Errorf("failed to connect to redis redis after %s", conf.ConnDeadline.String())
		case <-ticker.C:
			cl, err = rueidis.NewClient(rueidis.ClientOption{
				InitAddress: []string{conf.Host},
				Password:    conf.Password,
			})
			if err == nil {
				l.Infof("connected to redis on host: %v ", conf.Host)

				return &cl, nil
			}
			l.Errorf("failed to connect to redis, reconnecting...: %v", err)

		}
	}

}
