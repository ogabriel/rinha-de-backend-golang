#!/bin/sh

if [ "$1" = 'release' ]; then
    exec /app/rinha
elif [ "$1" = 'migrate_and_release' ]; then
    psql $DATABASE_URL -c "CREATE DATABASE $DATABASE_NAME"
    psql $DATABASE_URL -c "DROP DATABASE $DATABASE_NAME"
    /app/migrate -path /app/migrations/ -database $DATABASE_URL/$DATABASE_NAME?sslmode=disable -verbose up

    exec /app/rinha
fi
