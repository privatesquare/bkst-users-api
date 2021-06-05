#!/bin/sh

mysqlsh --host localhost -P 3306 -u bkst-admin -D bkstusersdb --sql < users.sql

for row in $(cat < users.json | jq '.[]' -c); do
  curl --location --request POST 'http://localhost:8080/users' \
  --header 'Content-Type: application/json' \
  --data-raw "${row}"
  printf "\n"
done