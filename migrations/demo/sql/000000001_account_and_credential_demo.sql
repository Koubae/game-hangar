-- +migrate Up
INSERT INTO account (username, email, id)
VALUES ('account_test_1', 'account_test_1@test.com', '06e1b677-a4fe-42cf-8afd-ceec867d1fa5'),
       ('account_test_2', 'account_test_2@test.com', '93a6437f-af81-47c1-9aa4-5caedc0a6869'),
       ('account_test_3', NULL, '9e0f95a2-6535-4679-9db6-c93a823702b1')
;

SELECT * FROM account;

-- Providers
INSERT INTO provider (id, name, display_name, category)
VALUES
    -- Managed
    (1, 'username',   'Username',            'managed'),
    (2, 'email',      'Email',               'managed'),
    (3, 'device',     'Device',              'managed'),
    -- Anonymous
    (4, 'anonymous',  'Anonymous',           'anonymous'),
    (5, 'guest',      'Guest',               'anonymous'),
    -- Platform
    (6, 'steam',      'Steam',               'platform'),
    (7, 'epic',       'Epic Games',          'platform'),
    (8, 'psn',        'PlayStation Network', 'platform'),
    (9, 'xbox',       'Xbox',                'platform'),
    (10, 'nintendo',   'Nintendo',            'platform'),
    (11, 'gpg',        'Google Play Games',   'platform'),
    (12, 'gamecenter', 'Apple Game Center',   'platform'),
    -- Social
    (13, 'google',     'Google',              'social'),
    (14, 'apple',      'Apple',               'social'),
    (15, 'discord',    'Discord',             'social'),
    (16, 'facebook',   'Facebook',            'social')
;


SELECT * FROM provider;


-- Credentials (using variables for account and provider IDs)
-- +migrate StatementBegin
DO
$$
DECLARE
    -- Account IDs (selected by email)
    account_1_id UUID := (SELECT id FROM account WHERE username = 'account_test_1');
    account_2_id UUID := (SELECT id FROM account WHERE username = 'account_test_2');
    account_3_id UUID := (SELECT id FROM account WHERE username = 'account_test_3');
    -- Provider IDs (selected by name)
    provider_username BIGINT := (SELECT id FROM provider WHERE name = 'username');
    provider_email    BIGINT := (SELECT id FROM provider WHERE name = 'email');
BEGIN
    INSERT INTO account_credentials (credential, account_id, provider_id, secret, verified, verified_at)
    VALUES
        -- user_test_1
        ('account_test_1', account_1_id, provider_username,
         encode(digest('pass', 'sha256'), 'hex'),
         true, CURRENT_TIMESTAMP),
        ('account_test_1@test.com', account_1_id, provider_email,
         encode(digest('pass', 'sha256'), 'hex'),
         true, CURRENT_TIMESTAMP),

        -- user_test_2
        ('account_test_2', account_2_id, provider_username,
         encode(digest('pass', 'sha256'), 'hex'),
         true, CURRENT_TIMESTAMP),
        ('account_test_2@test.com', account_2_id, provider_email,
         encode(digest('pass', 'sha256'), 'hex'),
         false, NULL),

        -- user_test_3
        ('account_test_3', account_3_id, provider_username,
         encode(digest('pass', 'sha256'), 'hex'),
         true, CURRENT_TIMESTAMP)
;
END
$$;
-- +migrate StatementEnd


SELECT * FROM account_credentials;


SELECT * FROM account WHERE id = '06e1b677-a4fe-42cf-8afd-ceec867d1fa5';
SELECT * FROM account_credentials WHERE account_id = '06e1b677-a4fe-42cf-8afd-ceec867d1fa5';


SELECT account.*,  provider.name, credentials.*
FROM account account
         JOIN account_credentials credentials ON account.id = credentials.account_id
         JOIN provider provider ON credentials.provider_id = provider.id
WHERE account.id = '06e1b677-a4fe-42cf-8afd-ceec867d1fa5'
  AND provider.disabled IS FALSE
  AND credentials.disabled IS FALSE
;

-- ///////////////////
-- // ROLLBACK
-- ///////////////////

-- +migrate Down
