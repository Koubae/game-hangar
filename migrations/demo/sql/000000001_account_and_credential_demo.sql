-- +migrate Up
INSERT INTO account (username, email, id)
VALUES ('account_test_1', 'account_test_1@test.com', '06e1b677-a4fe-42cf-8afd-ceec867d1fa5'),
       ('account_test_2', 'account_test_2@test.com', '93a6437f-af81-47c1-9aa4-5caedc0a6869'),
       ('account_test_3', NULL, '9e0f95a2-6535-4679-9db6-c93a823702b1')
;

SELECT * FROM account;

-- Providers
INSERT INTO provider (id, source, type, display_name, category)
VALUES
    -- Managed
    (1,  'global', 'username',   'Username',            'managed'),
    (2,  'global', 'email',      'Email',               'managed'),
    (3,  'global', 'device',     'Device',              'managed'),
    -- Anonymous
    (4,  'global', 'anonymous',  'Anonymous',           'anonymous'),
    (5,  'global', 'guest',      'Guest',               'anonymous'),
    -- Platform
    (6,  'global', 'steam',      'Steam',               'platform'),
    (7,  'global', 'epic',       'Epic Games',          'platform'),
    (8,  'global', 'psn',        'PlayStation Network', 'platform'),
    (9,  'global', 'xbox',       'Xbox',                'platform'),
    (10, 'global', 'nintendo',   'Nintendo',            'platform'),
    (11, 'global', 'gpg',        'Google Play Games',   'platform'),
    (12, 'global', 'gamecenter', 'Apple Game Center',   'platform'),
    -- Social
    (13, 'global', 'google',     'Google',              'social'),
    (14, 'global', 'apple',      'Apple',               'social'),
    (15, 'global', 'discord',    'Discord',             'social'),
    (16, 'global', 'facebook',   'Facebook',            'social')
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
    provider_username BIGINT := (SELECT id FROM provider WHERE source = 'global' AND type = 'username');
    provider_email    BIGINT := (SELECT id FROM provider WHERE source = 'global' AND type = 'email');
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

        -- NOTE: For, I prefer to keep 1 UNIQUE account => provider type because is easier to start
        -- be more "restrictive" and the drop the key rather than allows duplicate and enforce later
        -- why this could be important: we could allow multiple email into an account and set the "default"
        -- being the one attached to the account.email 
        -- ('account_test_1_secondary@test.com', account_1_id, provider_email,
        --  encode(digest('pass', 'sha256'), 'hex'),
        --  true, CURRENT_TIMESTAMP),

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


SELECT account.*, provider.source, provider.type, credentials.*
FROM account account
         JOIN account_credentials credentials ON account.id = credentials.account_id
         JOIN provider provider ON credentials.provider_id = provider.id
WHERE account.id = '06e1b677-a4fe-42cf-8afd-ceec867d1fa5'
  AND provider.disabled IS FALSE
  AND credentials.disabled IS FALSE
;


SELECT account.*, provider.source, provider.type, credentials.*
FROM account account
         JOIN account_credentials credentials ON account.id = credentials.account_id
         JOIN provider provider ON credentials.provider_id = provider.id
WHERE account.username = 'account_test_1'
  AND provider.disabled IS FALSE
  AND credentials.disabled IS FALSE
;

SELECT account.*, credentials.*
FROM account account
         JOIN account_credentials credentials ON account.id = credentials.account_id
WHERE account.username = 'account_test_1'
  AND credentials.disabled IS FALSE
;

-- NOTE: The below information is completly wrong but I leave it here to have a laugh
-- and anways still I want to keep these 2 queries 🤣
-- NOTE: 2 Actually I was right. Because UNIQUE (account_id, provider_id) then 
-- If you get anything from an account with specific provider than means implicity 
-- that credential is the ONLY one really possibly linked and there is no need to grab by credential too
-- My theory is that below 2 query are equivalent as long as key 
--  UNIQUE (provider_id, credential) exists 

SELECT * FROM account_credentials 
WHERE 1=1
    AND account_id = '06e1b677-a4fe-42cf-8afd-ceec867d1fa5'
    AND provider_id = 2 
    AND credential = 'account_test_1@test.com'
;

SELECT * FROM account_credentials 
WHERE 1=1
    AND account_id = '06e1b677-a4fe-42cf-8afd-ceec867d1fa5'
    AND provider_id = 2 
;
 


-- ///////////////////
-- // ROLLBACK
-- ///////////////////

-- +migrate Down
