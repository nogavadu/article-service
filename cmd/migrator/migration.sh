#!/bin/bash
export MIGRATION_DSN="host=$PG_HOST port=$PG_PORT dbname=$PG_DB_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable"

sleep 2 && goose -dir "migrations/" postgres "${MIGRATION_DSN}" up -v