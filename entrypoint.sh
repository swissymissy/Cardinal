#!/bin/sh
# exit immediately if a command fails
set -e

# run migration
echo "Running migration..."
until goose -dir /app/sql/schema postgres "$DB_URL" up; do 
    echo "Migration failed, retrying in 2 seconds..."
    sleep 2
done

# hand off to the main app
echo "Starting server..."
exec /app/cardinal