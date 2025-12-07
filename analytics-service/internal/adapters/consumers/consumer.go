package consumers

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/analytics/service"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/analytics-service/internal/models"
)

type AnalyticsConsumer struct {
	service *service.AnalyticsService
}

func NewAnalyticsConsumer(analyticsService *service.AnalyticsService) *AnalyticsConsumer {
	return &AnalyticsConsumer{
		service: analyticsService,
	}
}

func (cons *AnalyticsConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (cons *AnalyticsConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (cons *AnalyticsConsumer) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	l := log.FromContext(s.Context())

	for {
		select {
		case msg, ok := <-c.Messages():
			if !ok {
				l.Info("Consumer channel closed")
				return nil
			}

			l.Infof("Message received: topic=%s, partition=%d, offset=%d, key=%s",
				msg.Topic, msg.Partition, msg.Offset, string(msg.Key))

			// Key содержит userId
			userUID := string(msg.Key)
			if userUID == "" {
				l.Errorf("Empty user_uid in message key")
				s.MarkMessage(msg, "")
				s.Commit()
				continue
			}

			// Value содержит действие (обычно "update")
			var kafkaMsg models.KafkaAnalyticsMessage
			if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
				// Если не удалось распарсить JSON, пробуем использовать просто строку
				l.Infof("Could not parse message as JSON, using key as userUID: %v", err)
				kafkaMsg.UserUID = userUID
				kafkaMsg.Action = string(msg.Value)
			}

			// Если userUID не в теле, берем из ключа
			if kafkaMsg.UserUID == "" {
				kafkaMsg.UserUID = userUID
			}

			l.Infof("Processing analytics for user: %s, action: %s",
				kafkaMsg.UserUID, kafkaMsg.Action)

			// Обрабатываем событие аналитики
			if err := cons.service.ProcessAnalyticsEvent(s.Context(), kafkaMsg.UserUID); err != nil {
				l.Errorf("Failed to process analytics event: %v", err)
				// Все равно помечаем как обработанное, чтобы избежать повторной обработки
				s.MarkMessage(msg, "")
				s.Commit()
				continue
			}

			l.Infof("Successfully processed analytics for user: %s", kafkaMsg.UserUID)

			// Помечаем сообщение как обработанное
			s.MarkMessage(msg, "")
			s.Commit()

		case <-s.Context().Done():
			l.Info("Consumer context done")
			return nil
		}
	}
}

func StartConsuming(ctx context.Context, consumer sarama.ConsumerGroup, topics []string, handler sarama.ConsumerGroupHandler) {
	l := log.FromContext(ctx)

	go func() {
		for {
			if err := consumer.Consume(ctx, topics, handler); err != nil {
				l.Errorf("Kafka consumer error: %s", err)

				if ctx.Err() != nil {
					return
				}
			}
		}
	}()
}
