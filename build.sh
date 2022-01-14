docker build -t registry.gitlab.com/leetcyber/leet-app:latest .

docker login registry.gitlab.com -u patrickbreen -p glpat-m-sAKKWW4huxAywb_J9A
docker push registry.gitlab.com/leetcyber/leet-app

