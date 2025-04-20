#!/bin/bash

if [ $# -lt 1 ]; then
  echo 出力ファイルは必須です.
elif [ $# -gt 1 ]; then
  echo 不正な引数です.
else
  migrate create -ext sql -dir db/migrations -seq $1
fi
