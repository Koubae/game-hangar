package testobj

const (
	PassHash = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	ProviderSourceDefault  = "global"
	ProviderSourceDisabled = "global-disabled"
	ProviderIDDisabled1    = 17
	ProviderIDDisabled2    = 18
	ProviderIDDisabled3    = 19

	CredentialAccount01 = "account_test_1"

	SQLAccountDemoData = `
INSERT INTO account (username, email, id)
VALUES ('account_test_1', 'account_test_1@test.com', '06e1b677-a4fe-42cf-8afd-ceec867d1fa5'),
       ('account_test_2', 'account_test_2@test.com', '93a6437f-af81-47c1-9aa4-5caedc0a6869'),
       ('account_test_3', NULL, '9e0f95a2-6535-4679-9db6-c93a823702b1')
;
	`

	SQLCredentialsDemoData = `
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

	`
)
