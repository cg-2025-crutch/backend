package models

//go:generate easyjson -all $GOFILE

type KafkaNotificationMessage struct {
	UserUID      string       `json:"user_uid"`
	Notification Notification `json:"notification"`
}
