-- +goose Up
ALTER TABLE user_invitations
ADD COLUMN expires_at TIMESTAMP(0)
WITH
    TIME ZONE NOT NULL DEFAULT (NOW () + INTERVAL '3 days');

-- +goose Down
ALTER TABLE user_invitations
DROP COLUMN expires_at;