# Authentication using Auth0

Golang gateway using Auth0 server

## Developers

```shell
# Generate token for app
curl --request POST \
  --url https://dev-k18jl6aj.us.auth0.com/oauth/token \
  --header 'content-type: application/json' \
  --data '{"client_id":"qDiQmvZl0euArFBDQK0C4R8k5g0bMZqD","client_secret":"YpK1MzUTgyQsM8OyztZC8hnLnAl-Ae_YDHTuByrAxrR52lhs9-DHdwyTbu6Ig27A","audience":"https://ruben/","grant_type":"client_credentials"}'

# Store token in env var
export TOKEN=

# Connect to endpoint using authentication
curl --request GET \
  --url http://localhost:3000/api/private \
  --header 'authorization: Bearer $TOKEN'
```
