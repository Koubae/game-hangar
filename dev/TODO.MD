TODO (Game Hangar)
==================


Identity
========

### Auth API Endpoints

```bash
POST /v1/auth/register/username
POST /v1/auth/register/email

# 
POST /v1/auth/login/username
POST /v1/auth/login/email
# 3rd party -- providers
POST /v1/auth/login/steam
POST /v1/auth/login/playstation
POST /v1/auth/login/device
POST /v1/auth/login/guest

POST /v1/auth/refresh
```

### Account API Endpoints


```bash
### Account Management - Backoffice
# List Accounts
GET     /v1/accounts
# Create Account
POST    /v1/accounts
# Get Account
GET     /v1/accounts/{account_id}
# Update Account
PATCH   /v1/accounts/{account_id}
# Delete Account
DELETE   /v1/accounts/{account_id}

# Moderation | Account status
POST /v1/accounts/{account_id}/suspend
POST /v1/accounts/{account_id}/restore
POST /v1/accounts/{account_id}/ban
POST /v1/accounts/{account_id}/unban


# Account Management - (as User)
# Get My Account
GET   /v1/accounts/me
# Update My Account
PATCH /v1/accounts/me
# disable account  -- soft delete
POST /v1/accounts/me/deactivate

# Manage Devices - (as User)
GET /v1/accounts/me/devices
GET /v1/accounts/me/devices/{device_id}
DELETE /v1/accounts/me/devices/{device_id}

### Link / Unlink Identities
## TODO: Check this was only brain stormed a bit with GPT. but need real thought process before implementing.
GET    /v1/accounts/me/identities
POST   /v1/accounts/me/identities/steam
POST   /v1/accounts/me/identities/playstation
POST   /v1/accounts/me/identities/email
POST   /v1/accounts/me/identities/username
POST   /v1/accounts/me/identities/device

DELETE /v1/accounts/me/identities/steam
DELETE /v1/accounts/me/identities/playstation
DELETE /v1/accounts/me/identities/email
DELETE /v1/accounts/me/identities/username
DELETE /v1/accounts/me/identities/device
```



### Create account curl


```bash
curl --location 'http://localhost:10000/api/v1/auth/register/username' \
--header 'Content-Type: application/json' \
--data '{
  "username": "user-1",
  "password": "test"
}'
```

