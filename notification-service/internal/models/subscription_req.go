package models

//go:generate easyjson -all $GOFILE

type SubscriptionReq struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
	UserAgent string `json:"user_agent,omitempty"`
	DeviceID  string `json:"device_id,omitempty"`
}
