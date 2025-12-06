package models

import (
	"time"
)

type Transaction struct {
	ID              int64     `json:"id" db:"id"`
	UserUID         string    `json:"user_uid" db:"user_uid"`
	CategoryID      int32     `json:"category_id" db:"category_id"`
	Type            string    `json:"type" db:"type"` // income or expense
	Amount          float64   `json:"amount" db:"amount"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	Category        *Category `json:"category,omitempty" db:"-"` // Populated by join
}

type CreateTransactionInput struct {
	UserUID         string    `json:"user_uid"`
	CategoryID      int32     `json:"category_id"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
}

type UpdateTransactionInput struct {
	ID              int64     `json:"id"`
	UserUID         string    `json:"user_uid"`
	CategoryID      int32     `json:"category_id"`
	Type            string    `json:"type"`
	Amount          float64   `json:"amount"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
}
