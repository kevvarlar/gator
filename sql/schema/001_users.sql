-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;