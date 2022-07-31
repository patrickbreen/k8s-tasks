### How to operate this API
```
#################################################################
# If you just want to run a dev version and don't want to touch kubernetes for the DB:
# Note you still need keycloak working (on kubernetes)
#################################################################

export POSTGRES_CONNECTION=dev
go build .
go test . -v
./leet

#################################################################
# Everything below assumes some amount of kubernetes usage:
#################################################################

# portforward postgres from the cluster to use in localhost dev/testing
k port-forward -n tasks svc/tasks-postgres-master 5432:5432

# get the password from a secret generated by postgres operator:
# kubectl get secret owner.tasks-postgres.credentials.postgresql.acid.zalan.do -n tasks -o 'jsonpath={.data.password}' | base64 -d

export POSTGRES_CONNECTION="host=localhost port=5432 user=owner dbname=app password=<password> sslmode=disable"

# connect with psql:
psql -h localhost -p 5432 -U owner -d app

# install dependencies:
go get .

# build and run app
go build
./leet

# build  and run canary
cd canary
go build
./canary

# build container images
docker build -f <docker file> -t <tag> .
```

