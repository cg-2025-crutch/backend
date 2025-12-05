package models

//go:generate easyjson -all $GOFILE

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Data  string `json:"data,omitempty"`
}
