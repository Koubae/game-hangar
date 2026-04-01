


SELECT id FROM account WHERE username = 'manager';
SELECT id FROM admin_account WHERE account_id = (SELECT id FROM account WHERE username = 'manager');
SELECT * FROM admin_account WHERE account_id = (SELECT id FROM account WHERE username = 'manager');

SELECT * FROM admin_account admin JOIN account_permissions permissions ON admin.id = permissions.admin_account_id WHERE permissions.permission_id = 1;

-- load admin account + permissions
SELECT * FROM admin_account admin JOIN account_permissions permissions ON admin.id = permissions.admin_account_id WHERE account_id = (SELECT id FROM account WHERE username = 'manager');
SELECT perm.id, perm.service, perm.resource, perm.action
    FROM admin_account admin
        JOIN account_permissions grants ON admin.id = grants.admin_account_id
        JOIN permissions perm ON grants.permission_id = perm.id
    WHERE account_id = (SELECT id FROM account WHERE username = 'manager');



SELECT perm.*
FROM admin_account admin
         JOIN account_permissions grants ON admin.id = grants.admin_account_id
         JOIN permissions perm ON grants.permission_id = perm.id
WHERE account_id = (SELECT id FROM account WHERE username = 'manager');

