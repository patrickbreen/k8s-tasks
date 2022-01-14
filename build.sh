
# format source code:
gofmt -w -s .

# build docker image
docker build -t registry.gitlab.com/leetcyber/leet-app:latest .

# push docker image
docker login registry.gitlab.com -u patrickbreen -p glpat-m-sAKKWW4huxAywb_J9A
docker push registry.gitlab.com/leetcyber/leet-app:latest
