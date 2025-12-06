package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *FundsRepository) updateUserBalanceInTx(ctx context.Context, tx pgx.Tx, userUID, transactionType string, amount float64, add bool) error {
	insertQuery := `
		INSERT INTO user_balances (user_uid, total_balance, total_income, total_expense, last_transaction_at, updated_at)
		VALUES ($1, 0, 0, 0, NOW(), NOW())
		ON CONFLICT (user_uid) DO NOTHING
	`
	_, err := tx.Exec(ctx, insertQuery, userUID)
	if err != nil {
		return fmt.Errorf("failed to ensure balance record: %w", err)
	}

	var updateQuery string
	if transactionType == "income" {
		if add {
			updateQuery = `
				UPDATE user_balances 
				SET total_income = total_income + $1,
				    total_balance = total_balance + $1,
				    last_transaction_at = NOW(),
				    updated_at = NOW()
				WHERE user_uid = $2
			`
		} else {
			updateQuery = `
				UPDATE user_balances 
				SET total_income = total_income - $1,
				    total_balance = total_balance - $1,
				    updated_at = NOW()
				WHERE user_uid = $2
			`
		}
	} else {
		if add {
			updateQuery = `
				UPDATE user_balances 
				SET total_expense = total_expense + $1,
				    total_balance = total_balance - $1,
				    last_transaction_at = NOW(),
				    updated_at = NOW()
				WHERE user_uid = $2
			`
		} else {
			updateQuery = `
				UPDATE user_balances 
				SET total_expense = total_expense - $1,
				    total_balance = total_balance + $1,
				    updated_at = NOW()
				WHERE user_uid = $2
			`
		}
	}

	_, err = tx.Exec(ctx, updateQuery, amount, userUID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}
