-- +goose Up
CREATE TABLE IF NOT EXISTS followers (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, follower_id)
);

-- +goose Down
DROP TABLE IF EXISTS followers;