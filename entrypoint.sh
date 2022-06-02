#!/bin/bash

set -e

until  nc -z -v -w30 db 5432
  do
    echo "Waiting for database connection..."
    sleep 1
  done
echo "Database available, starting application"

go run cmd/main.go
