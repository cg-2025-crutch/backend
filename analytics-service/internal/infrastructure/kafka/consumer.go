package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/config"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/infrastructure/log"
	"github.com/google/uuid"
)

const (
	ReconnectPeriod = 5 * time.Second
)

func InitKafkaConsumer(ctx context.Context, appName string, conf config.KafkaConfig) (sarama.ConsumerGroup, error) {
	l := log.FromContext(ctx)

	clientUUID := uuid.New().String()
	clientID := fmt.Sprintf("%s-%s", appName, clientUUID)

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_3_0_0
	cfg.ClientID = clientID
	cfg.Consumer.Group.InstanceId = clientUUID
	cfg.Consumer.Offsets.AutoCommit.Enable = false
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	ticker := time.NewTicker(ReconnectPeriod)
	ctxStop, cancel := context.WithTimeout(ctx, conf.ConnDeadline)
	defer cancel()
	defer ticker.Stop()

	for {
		select {
		case <-ctxStop.Done():
			return nil, fmt.Errorf("failed to connect to consumer after %s", conf.ConnDeadline.String())
		case <-ticker.C:
			consumer, err := sarama.NewConsumerGroup(conf.Brokers, appName, cfg)
			if err == nil {
				l.Infof("Kafka consumer created with ID '%s'", clientID)
				return consumer, nil
			}
			l.Infof("Failed to create Kafka consumer group, retrying...:%v", err)
		}
	}
}
