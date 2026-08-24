-- +goose Up
CREATE TABLE
    IF NOT EXISTS comments (
        id BIGSERIAL PRIMARY KEY,
        post_id BIGINT REFERENCES posts (id) ON DELETE CASCADE,
        user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- +goose Down
DROP TABLE IF EXISTS comments;