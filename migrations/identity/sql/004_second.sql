-- +migrate Up
CREATE TABLE IF NOT EXISTS example_3 (
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- +migrate Down
DROP TABLE IF EXISTS example_3;
