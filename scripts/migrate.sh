#!/bin/bash

if [ $# != 1 ] || [ $1 != "up" ] && [ $1 != "down" ]; then
  echo 不正な引数です
else
  migrate -path db/migrations -database "mysql://$DATABASE_USER:$DATABASE_PASSWORD@tcp($DATABASE_HOST:$DATABASE_PORT)/$DATABASE_NAME" $1
fi
