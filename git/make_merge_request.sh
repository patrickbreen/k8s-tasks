set -x
cd ..
git clone "https://oath2:$MAKE_MERGE_REQUEST_TOKEN@gitlab.com/leetcyber/demo-k8s-infra.git"
cd "demo-k8s-infra"
git checkout  -b "update_container_image_tag_$CI_COMMIT_SHA"
sed -i "s/tag: \".*\"/tag: \"$CI_COMMIT_SHA\"/g" ./leet/gamma-values.yaml
sed -i "s/tag: \".*\"/tag: \"$CI_COMMIT_SHA\"/g" ./leet/prod-values.yaml
./build_helm.sh
cat ./leet/gamma-values.yaml
cat ./leet/gamma/manifest.yaml
git config --global user.email "merge_request_api@example.com"
git config --global user.name "merge request api"
git commit -am "updated container image tag"
git push \
  -o merge_request.create \
  -o merge_request.target=main \
  origin update_container_image_tag_$CI_COMMIT_SHA \
