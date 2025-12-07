package models

// KafkaAnalyticsMessage представляет сообщение из Kafka для аналитики
type KafkaAnalyticsMessage struct {
	UserUID string `json:"user_uid"`
	Action  string `json:"action"` // "update"
}
