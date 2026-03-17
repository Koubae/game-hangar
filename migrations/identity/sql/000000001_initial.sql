-- +migrate Up

-- ///////////////////
-- // Extensions
-- ///////////////////
CREATE EXTENSION IF NOT EXISTS pgcrypto;   -- gen_random_uuid(), digest(), etc.
CREATE EXTENSION IF NOT EXISTS citext;     -- case-insensitive text (perfect for email)

-- ///////////////////
-- // Schemas
-- ///////////////////
-- I decided to not create the schema here for now since I belive it complicates things too much.
-- However It may be useful in the future and I just leave it here as reference for future me or future anyone.
-- CREATE SCHEMA IF NOT EXISTS identity;


-- ///////////////////
-- // Triggers
-- ///////////////////
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_current_timestamp()
    RETURNS trigger
    LANGUAGE plpgsql AS
$$
BEGIN
    IF NEW IS DISTINCT FROM OLD THEN
        NEW.updated := CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END
$$;
-- +migrate StatementEnd

-- ///////////////////
-- // ROLLBACK
-- ///////////////////

-- +migrate Down

DROP FUNCTION IF EXISTS update_current_timestamp();

DROP EXTENSION IF EXISTS citext;
DROP EXTENSION IF EXISTS pgcrypto;
