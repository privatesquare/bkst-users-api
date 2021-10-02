#!/bin/zsh
orgName="privatesquare"
appName="bkst-users-api"
outputPath="docker-context"

docker compose -f docker-context/docker-compose.yml -p "${appName}"  up -d

 sleep 30

 mysqlsh --host localhost -P 3306 -u bkst-admin -D bkstusersdb --sql < resources/users.sql

 go fmt ./...
 go build ./...

 ./"${appName}"

for row in $(cat < resources/users.json | jq '.[]' -c); do
  curl --location --request POST 'http://localhost:8080/users' \
  --header 'Content-Type: application/json' \
  --data-raw "${row}"
  printf "\n"
done
