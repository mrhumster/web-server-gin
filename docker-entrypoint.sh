#!/bin/sh
set -e
for i in $(seq 1 30); do
  pg_isready "$DB_HOST" && break || sleep 2
done
for i in $(seq 1 5); do
  migrate -path /app/migrations -database "$DATABASE_URL" up && break || sleep $((i * 5))
done
