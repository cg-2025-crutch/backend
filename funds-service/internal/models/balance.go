package models

import (
	"time"
)

type UserBalance struct {
	UserUID           string     `json:"user_uid" db:"user_uid"`
	TotalBalance      float64    `json:"total_balance" db:"total_balance"`
	TotalIncome       float64    `json:"total_income" db:"total_income"`
	TotalExpense      float64    `json:"total_expense" db:"total_expense"`
	LastTransactionAt *time.Time `json:"last_transaction_at,omitempty" db:"last_transaction_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}
