package consumers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/cg-2025-crutch/backend/notification-service/internal/infrastructure/log"
	"github.com/cg-2025-crutch/backend/notification-service/internal/models"
	"github.com/cg-2025-crutch/backend/notification-service/internal/notifications/service"
	"github.com/mailru/easyjson"
)

type NotificationConsumer struct {
	service *service.NotificationService
}

func NewNotificationConsumer(notifService *service.NotificationService) *NotificationConsumer {
	return &NotificationConsumer{
		service: notifService,
	}
}

func (cons *NotificationConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (cons *NotificationConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (cons *NotificationConsumer) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	l := log.FromContext(s.Context())

	for {
		select {
		case msg, ok := <-c.Messages():
			if !ok {
				l.Info("Consumer channel closed")
				return nil
			}

			l.Infof("Message received: topic=%s, partition=%d, offset=%d",
				msg.Topic, msg.Partition, msg.Offset)

			var kafkaMsg models.KafkaNotificationMessage
			if err := easyjson.Unmarshal(msg.Value, &kafkaMsg); err != nil {
				l.Errorf("Failed to unmarshal Kafka message: %v", err)
				// Mark message as processed even if it fails to avoid blocking
				s.MarkMessage(msg, "")
				s.Commit()
				continue
			}

			l.Infof("Processing notification for user: %s, title: %s",
				kafkaMsg.UserUID, kafkaMsg.Notification.Title)

			if err := cons.service.SendNotification(s.Context(), kafkaMsg.UserUID, kafkaMsg.Notification); err != nil {
				l.Errorf("Failed to send notification: %v", err)
				s.MarkMessage(msg, "")
				s.Commit()
				continue
			}

			l.Infof("Successfully sent notification to user: %s", kafkaMsg.UserUID)

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
