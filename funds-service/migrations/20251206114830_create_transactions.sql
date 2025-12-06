-- +goose Up
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_uid UUID NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories(id),
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    transaction_date DATE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS transactions;
