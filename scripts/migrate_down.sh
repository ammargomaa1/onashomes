#!/bin/bash

# Script to rollback migrations
# Usage: ./scripts/migrate_down.sh [number_of_migrations]

MIGRATION_DIR="migrations"
STEPS=${1:-1}

echo "⚠️  Rolling back $STEPS migration(s)..."

# Find the last N down migration files
DOWN_FILES=$(ls -r ${MIGRATION_DIR}/*.down.sql 2>/dev/null | head -n $STEPS)

if [ -z "$DOWN_FILES" ]; then
    echo "No migration files found in ${MIGRATION_DIR}"
    exit 1
fi

# Load database configuration from .env
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if required env variables are set
if [ -z "$DB_HOST" ] || [ -z "$DB_NAME" ] || [ -z "$DB_USER" ]; then
    echo "Error: Database configuration not found in .env file"
    exit 1
fi

# Execute each down migration
for file in $DOWN_FILES; do
    echo "Executing: $file"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $file
    
    if [ $? -eq 0 ]; then
        echo "✓ Successfully executed: $file"
    else
        echo "✗ Failed to execute: $file"
        exit 1
    fi
done

echo "✓ Rollback completed successfully"
