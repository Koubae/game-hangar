-- +migrate Up
INSERT INTO account (username, email, id)
VALUES
    ('root', NULL, 'c7b9d2c4-8f1a-4e3b-9c5e-2a6f4d1a0b7c'),
    ('admin', NULL, '8f3a1c6e-5d4b-4f2f-9c21-7e0d3a9b6c14'),
    ('manager', NULL, '2d9b4f51-6c87-4d9e-a8b3-1f2c0d7e9a4b'),
    ('employee', NULL, 'a1f0c3e9-7b4d-4d21-9c6a-5e8f2b3d1a70')
;




INSERT INTO permissions (service, resource, action)
VALUES
    -- Superuser
    ( '*', '*', '*'),
    -- ==================== < Identity > ==================== --
    ( 'identity', '*', '*'),
    -- Auth
    ('identity', 'auth', '*'),
    ('identity', 'auth', 'read'),
    ('identity', 'auth', 'write'),
    ('identity', 'auth', 'delete'),
    -- Account
    ('identity', 'account', '*'),
    ('identity', 'account', 'read'),
    ('identity', 'account', 'write'),
    ('identity', 'account', 'delete')
;

-- +migrate StatementBegin
DO
$$
    DECLARE
        -- Account IDs (selected by email)
        root UUID := (SELECT id FROM account WHERE username = 'root');
        admin UUID := (SELECT id FROM account WHERE username = 'admin');
        manager UUID := (SELECT id FROM account WHERE username = 'manager');
        employee UUID := (SELECT id FROM account WHERE username = 'employee');
    BEGIN
        INSERT INTO admin_account (account_id)
        VALUES
            -- root
            (root),
            (admin),
            (manager),
            (employee)
;

    END
$$;
-- +migrate StatementEnd

-- Credentials + Admin Account + Permissions Scopes (using variables for account and provider IDs)
-- +migrate StatementBegin
DO
$$
DECLARE
    -- Provider IDs (selected by name)
    provider_username BIGINT := (SELECT id FROM provider WHERE source = 'global' AND type = 'username');

    -- Account IDs (selected by email)
    root UUID := (SELECT id FROM account WHERE username = 'root');
    admin_root BIGINT := (SELECT id FROM admin_account WHERE account_id = root);

    admin UUID := (SELECT id FROM account WHERE username = 'admin');
    admin_admin BIGINT := (SELECT id FROM admin_account WHERE account_id = admin);

    manager UUID := (SELECT id FROM account WHERE username = 'manager');
    admin_manager BIGINT := (SELECT id FROM admin_account WHERE account_id = manager);

    employee UUID := (SELECT id FROM account WHERE username = 'employee');
    admin_employee BIGINT := (SELECT id FROM admin_account WHERE account_id = employee);

    permission_su BIGINT := (SELECT id FROM permissions WHERE service = '*' AND resource = '*' AND action = '*');

    permission_identity_su BIGINT := (SELECT id FROM permissions WHERE service = 'identity' AND resource = '*' AND action = '*');

    permission_identity_auth_su BIGINT := (SELECT id FROM permissions WHERE service = 'identity' AND resource = 'auth' AND action = '*');
    permission_identity_auth_read BIGINT := (SELECT id FROM permissions WHERE service = 'identity' AND resource = 'auth' AND action = 'read');
    permission_identity_auth_write BIGINT := (SELECT id FROM permissions WHERE service = 'identity' AND resource = 'auth' AND action = 'write');
    permission_identity_auth_delete BIGINT := (SELECT id FROM permissions WHERE service = 'identity' AND resource = 'auth' AND action = 'delete');

    permission_identity_account_su BIGINT := (SELECT id FROM permissions WHERE service = 'identity' AND resource = 'account' AND action = '*');
    permission_identity_account_read BIGINT := (SELECT id FROM permissions WHERE service = 'identity' AND resource = 'account' AND action = 'read');
    permission_identity_account_write BIGINT := (SELECT id FROM permissions WHERE service = 'identity' AND resource = 'account' AND action = 'write');
    permission_identity_account_delete BIGINT := (SELECT id FROM permissions WHERE service = 'identity' AND resource = 'account' AND action = 'delete');

BEGIN
    INSERT INTO account_credentials (credential, account_id, provider_id, secret, verified, verified_at)
    VALUES
        -- root
        ('root', root, provider_username,
         '$2a$10$/ibRSFWvBxGUO2lYMnh0yOTGSIycp.Ae6Oc5Py/fIVYBxpT5PTGAS',
         true, CURRENT_TIMESTAMP),
        ('admin', admin, provider_username,
         '$2a$10$/ibRSFWvBxGUO2lYMnh0yOTGSIycp.Ae6Oc5Py/fIVYBxpT5PTGAS',
         true, CURRENT_TIMESTAMP),
        ('manager', manager, provider_username,
         '$2a$10$/ibRSFWvBxGUO2lYMnh0yOTGSIycp.Ae6Oc5Py/fIVYBxpT5PTGAS',
         true, CURRENT_TIMESTAMP),
        ('employee', employee, provider_username,
         '$2a$10$/ibRSFWvBxGUO2lYMnh0yOTGSIycp.Ae6Oc5Py/fIVYBxpT5PTGAS',
         true, CURRENT_TIMESTAMP)
;

    INSERT INTO account_permissions (admin_account_id, permission_id)
    VALUES
        (admin_root, permission_su),

        (admin_admin, permission_identity_su),

        (admin_manager, permission_identity_auth_read),
        (admin_manager, permission_identity_auth_write),
        (admin_manager, permission_identity_account_read),
        (admin_manager, permission_identity_account_write),

        (admin_employee, permission_identity_auth_read),
        (admin_employee, permission_identity_account_read)
;
END
$$;
-- +migrate StatementEnd


-- ///////////////////
-- // ROLLBACK
-- ///////////////////

-- +migrate Down
