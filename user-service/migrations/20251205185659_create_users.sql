-- +goose Up
CREATE TABLE users(
    uid UUID PRIMARY KEY,
    username TEXT,
    password TEXT,
    first_name TEXT,
    second_name TEXT,
    age INT,
    salary DOUBLE PRECISION,
    work_sphere INT

);

-- +goose Down
DROP TABLE IF EXISTS users;
