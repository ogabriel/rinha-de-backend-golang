#!/bin/sh

if [ "$1" = 'release' ]; then
    exec /app/rinha-de-backend-golang
elif [ "$1" = 'migrate_and_release' ]; then
    make database-reset
    exec /app/rinha-de-backend-golang
fi
