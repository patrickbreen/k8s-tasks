docs for api key example:
https://keycloak.discourse.group/t/how-to-use-api-keys-with-keycloak/2390/13

```
#build:
go build .

#run:
sudo ./oauth

#test:
go to "oauth.dev.leetcyber.com" which directs to localhost:80 in /etc/hosts

#get api key: curl -v -k -d "client_id=client-secret" -d "username=patrick" -d "password=star" -d "grant_type=password" -X POST "https://keycloak.dev.leetcyber.com/auth/realms/basic/protocol/openid-connect/token"
```

TODO:

### make an http client:

- get token

- loop:
  - if token expired, refresh token,
  - else make http request to server
  - verify that server returns expected response

- make tasks service check for, and validate bearer token in middleware
- take mTLS off of tasks service
- make client requests tokens and refreshes, and send bearer token
