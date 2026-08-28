-- +goose Up
CREATE TABLE
    IF NOT EXISTS roles (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        level INT NOT NULL DEFAULT 0,
        description TEXT
    );

INSERT INTO
    roles (name, level, description)
VALUES
    (
        'user',
        1,
        'a user can only access their own data'
    ),
    (
        'moderator',
        2,
        'a moderator can update other users posts'
    ),
    (
        'admin',
        3,
        'an admin can update and delete other users posts'
    );

-- +goose Down
DROP TABLE IF EXISTS roles;