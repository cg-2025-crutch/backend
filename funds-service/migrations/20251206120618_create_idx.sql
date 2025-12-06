-- +goose Up
CREATE INDEX idx_transactions_user_date ON transactions(user_uid, transaction_date DESC);
CREATE INDEX idx_transactions_user_category ON transactions(user_uid, category_id);
CREATE INDEX idx_transactions_user_type_date ON transactions(user_uid, type, transaction_date);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);

-- +goose Down
ALTER TABLE transactions DROP CONSTRAINT idx_transactions_user_date;
ALTER TABLE transactions DROP CONSTRAINT idx_transactions_user_category;
ALTER TABLE transactions DROP CONSTRAINT idx_transactions_user_type_date;
ALTER TABLE transactions DROP CONSTRAINT idx_transactions_date;
