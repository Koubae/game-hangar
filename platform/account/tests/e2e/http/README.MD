Testing e2e -- http
===================

```bash
access_token=$(curl -s --request POST  localhost:8001/api/v1/auth/login -H "Content-Type: application/json" -d '{"username": "user_1","password": "pass","client_id": "game_1.x.y.z"}' | jq -r '.access_token'); \
  echo $access_token; \
  curl localhost:8001/api/v1/account/user_1?full_profile=false -H "Authorization: Bearer $access_token"
  
```