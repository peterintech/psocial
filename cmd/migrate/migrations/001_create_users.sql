-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE
    IF NOT EXISTS users (
        id BIGSERIAL PRIMARY KEY,
        email citext UNIQUE NOT NULL,
        username VARCHAR(255) UNIQUE NOT NULL,
        password bytea NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- +goose Down
DROP TABLE IF EXISTS users;