package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	producer "github.com/cg-2025-crutch/backend/funds-service/internal/adapters/producers"
	"github.com/cg-2025-crutch/backend/funds-service/internal/config"
	"github.com/cg-2025-crutch/backend/funds-service/internal/infrastructure/log"
	"github.com/google/uuid"
)

func InitKafkaProducer(ctx context.Context, appName string, conf config.KafkaConfig) (producer.Producer, error) {
	l := log.FromContext(ctx)

	clientUUID := uuid.New().String()
	clientID := fmt.Sprintf("%s-%s", appName, clientUUID)

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_3_0_0
	cfg.ClientID = clientID
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	ticker := time.NewTicker(10 * time.Second)
	ctxStop, cancel := context.WithTimeout(ctx, conf.ConnDeadline)
	defer cancel()
	defer ticker.Stop()

	for {
		select {
		case <-ctxStop.Done():
			return nil, fmt.Errorf("failed to connect to producer after %s", conf.ConnDeadline.String())
		case <-ticker.C:
			syncProducer, err := sarama.NewSyncProducer(conf.Brokers, cfg)

			if err == nil {
				producerAdapter := producer.NewKafkaProducer(syncProducer, conf.ProdTopic)
				return producerAdapter, nil
			}
			l.Infof("Failed to create Kafka producer, retrying...:%v", err)
		}
	}

}
