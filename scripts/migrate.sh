#!/bin/bash

# Migration management script
# Usage: 
#   ./scripts/migrate.sh up     - Run pending migrations
#   ./scripts/migrate.sh down   - Rollback last migration
#   ./scripts/migrate.sh status - Show migration status

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

case "$1" in
    "up")
        echo "🔄 Running migrations..."
        cd "$PROJECT_ROOT"
        go run cmd/api/main.go --migrate-only
        ;;
    "down")
        echo "🔄 Rolling back last migration..."
        cd "$PROJECT_ROOT"
        go run cmd/api/main.go --rollback
        ;;
    "status")
        echo "📋 Migration status:"
        cd "$PROJECT_ROOT"
        go run cmd/api/main.go --migration-status
        ;;
    "create")
        if [ -z "$2" ]; then
            echo "Error: Migration name is required"
            echo "Usage: ./scripts/migrate.sh create migration_name"
            exit 1
        fi
        ./scripts/create_migration.sh "$2"
        ;;
    *)
        echo "Usage: $0 {up|down|status|create <name>}"
        echo ""
        echo "Commands:"
        echo "  up              Run all pending migrations"
        echo "  down            Rollback the last migration"
        echo "  status          Show migration status"
        echo "  create <name>   Create a new migration file"
        echo ""
        echo "Examples:"
        echo "  $0 up"
        echo "  $0 create add_user_table"
        echo "  $0 down"
        exit 1
        ;;
esac
