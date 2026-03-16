-- +migrate Up
CREATE TABLE IF NOT EXISTS account (
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- +migrate Down
DROP TABLE IF EXISTS account;
