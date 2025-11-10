#!/bin/bash

# Script to create a new migration file
# Usage: ./scripts/create_migration.sh migration_name

if [ -z "$1" ]; then
    echo "Error: Migration name is required"
    echo "Usage: ./scripts/create_migration.sh migration_name"
    exit 1
fi

MIGRATION_NAME=$1
TIMESTAMP=$(date +%Y%m%d%H%M%S)
MIGRATION_DIR="migrations"

# Create migrations directory if it doesn't exist
mkdir -p $MIGRATION_DIR

# Create up migration file
UP_FILE="${MIGRATION_DIR}/${TIMESTAMP}_${MIGRATION_NAME}.up.sql"
DOWN_FILE="${MIGRATION_DIR}/${TIMESTAMP}_${MIGRATION_NAME}.down.sql"

# Create empty migration files
touch $UP_FILE
touch $DOWN_FILE

echo "-- Migration: ${MIGRATION_NAME}" > $UP_FILE
echo "-- Created at: $(date)" >> $UP_FILE
echo "" >> $UP_FILE
echo "-- Add your UP migration SQL here" >> $UP_FILE

echo "-- Migration: ${MIGRATION_NAME}" > $DOWN_FILE
echo "-- Created at: $(date)" >> $DOWN_FILE
echo "" >> $DOWN_FILE
echo "-- Add your DOWN migration SQL here" >> $DOWN_FILE

echo "✓ Migration files created:"
echo "  - $UP_FILE"
echo "  - $DOWN_FILE"
