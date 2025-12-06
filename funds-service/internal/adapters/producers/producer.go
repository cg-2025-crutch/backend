package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/cg-2025-crutch/backend/funds-service/internal/infrastructure/log"
	"github.com/google/uuid"
)

type Producer interface {
	Produce(ctx context.Context, key, message []byte) error
	ProduceWithID(ctx context.Context, key, message, msgID []byte) error
	ProduceWithHeaders(ctx context.Context, key, message []byte, headers map[string]string) error
	Close() error
}

const msgIdHeader = "msg-id"

type notificationSender struct {
	topic    string
	producer sarama.SyncProducer
}

func NewKafkaProducer(producer sarama.SyncProducer, topic string) Producer {
	return &notificationSender{
		topic:    topic,
		producer: producer,
	}
}

func (p notificationSender) Produce(ctx context.Context, key, message []byte) error {
	msgUUID := uuid.New().String()

	return p.produce(ctx, []byte(msgUUID), key, message, nil)
}

func (p notificationSender) ProduceWithID(ctx context.Context, key, message, msgID []byte) error {
	return p.produce(ctx, msgID, key, message, nil)
}

func (p notificationSender) ProduceWithHeaders(ctx context.Context, key, message []byte, headers map[string]string) error {
	msgUUID := uuid.New().String()

	return p.produce(ctx, []byte(msgUUID), key, message, headers)
}

func (p notificationSender) produce(ctx context.Context, msgID, key, message []byte, headers map[string]string) error {
	l := log.FromContext(ctx)

	sHeaders := []sarama.RecordHeader{
		{
			Key:   []byte(msgIdHeader),
			Value: msgID,
		},
	}

	for k, v := range headers {
		sHeaders = append(sHeaders, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		l.Errorf("Failed to send message to Kafka topic '%s': %v", p.topic, err)
		return fmt.Errorf("failed to send message to Kafka topic '%s': %v", p.topic, err)
	}

	l.Infof("Message succesfully sent, partition = %d, offset = %d", partition, offset)
	return nil
}

func (n notificationSender) Close() error {
	//TODO implement me
	panic("implement me")
}
