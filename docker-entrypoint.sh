#!/bin/sh

if [ "$1" = 'release' ]; then
    make database-check
    exec /app/rinha-de-backend-golang
elif [ "$1" = 'migrate_and_release' ]; then
    make database-reset-release
    exec /app/rinha-de-backend-golang
fi
