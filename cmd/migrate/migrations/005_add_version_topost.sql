-- +goose Up
ALTER TABLE posts
ADD COLUMN IF NOT EXISTS version INT DEFAULT 1;

-- +goose Down
ALTER TABLE posts
DROP COLUMN IF EXISTS version;