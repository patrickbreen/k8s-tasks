# portforward postgres from the cluster to use in localhost dev/testing
k port-forward -n postgres service/my-bitnami-postgres-postgresql 5432:5432

# good example
https://github.com/gothinkster/golang-gin-realworld-example-app
