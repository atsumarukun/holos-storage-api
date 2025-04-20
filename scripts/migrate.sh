#!/bin/bash

if [ $# != 1 ] || [ $1 != "up" ] && [ $1 != "down" ]; then
  echo 不正な引数です
else
  migrate -path db/migrations -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" $1
fi
