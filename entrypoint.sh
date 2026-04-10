#!/bin/sh
# exit immediately if a command fails
set -e

# 1. wait for database
echo "Waiting for database..."
until nc -z $DB_HOST $DB_PORT; do 
    sleep 1
done 

# 2. run migration
echo "Running migration..."
goose -dir /app/sql/schema postgres "$DB_URL" up

# hand off to the main app
echo "Starting server..."
exec /app/cardinal