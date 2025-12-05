package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  int64     `json:"iat"`
	ExpiresAt int64     `json:"exp"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
