-- +goose Up
CREATE TABLE user_balances (
    user_uid UUID PRIMARY KEY,
    total_balance DECIMAL(15, 2) DEFAULT 0,
    total_income DECIMAL(15, 2) DEFAULT 0,
    total_expense DECIMAL(15, 2) DEFAULT 0,
    last_transaction_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS user_balances;