# Database Migrations

This document explains how to use the SQL migration system in the OnasHomes API.

## Overview

The application now includes a complete SQL migration system that:
- ✅ **Scans** the `migrations/` directory for `.sql` files
- ✅ **Executes** pending migrations in order
- ✅ **Tracks** applied migrations in the database
- ✅ **Supports** rollbacks with `.down.sql` files
- ✅ **Runs automatically** on application startup

## Migration File Naming Convention

Migration files follow this pattern (from `scripts/create_migration.sh`):
```
{TIMESTAMP}_{MIGRATION_NAME}.up.sql    # Forward migration
{TIMESTAMP}_{MIGRATION_NAME}.down.sql  # Rollback migration
```

Example:
```
20241112143022_add_user_table.up.sql
20241112143022_add_user_table.down.sql
```

## Creating Migrations

### Method 1: Using the Script
```bash
./scripts/create_migration.sh add_user_table
```

This creates:
- `migrations/20241112143022_add_user_table.up.sql`
- `migrations/20241112143022_add_user_table.down.sql`

### Method 2: Manual Creation
Create files manually following the naming convention.

## Migration Management

### Run All Pending Migrations
```bash
# Option 1: Start the full application (runs migrations automatically)
go run cmd/api/main.go

# Option 2: Run migrations only
go run cmd/api/main.go --migrate-only

# Option 3: Using the helper script
./scripts/migrate.sh up
```

### Rollback Last Migration
```bash
# Option 1: Direct command
go run cmd/api/main.go --rollback

# Option 2: Using the helper script
./scripts/migrate.sh down
```

### Check Migration Status
```bash
# Option 1: Direct command
go run cmd/api/main.go --migration-status

# Option 2: Using the helper script
./scripts/migrate.sh status
```

### Create New Migration
```bash
./scripts/migrate.sh create migration_name
```

## Migration Execution Order

1. **SQL Migrations** (`.up.sql` files) - Run first
2. **GORM AutoMigrate** - Runs after SQL migrations
3. **Seed Default Data** - Creates roles and admin user
4. **Permission Scanning** - Discovers and syncs permissions

## Example Migration Files

### Forward Migration (`20241112143022_add_products.up.sql`)
```sql
-- Migration: add_products
-- Created at: 2024-11-12

CREATE TABLE products (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);
```

### Rollback Migration (`20241112143022_add_products.down.sql`)
```sql
-- Migration: add_products rollback
-- Created at: 2024-11-12

DROP INDEX IF EXISTS idx_products_name;
DROP TABLE IF EXISTS products;
```

## Migration Tracking

The system automatically creates a `migration_records` table to track:
- ✅ Which migrations have been applied
- ✅ When they were applied
- ✅ Order of execution

## Features

### Automatic Execution
- Migrations run automatically when you start the application
- Only pending migrations are executed
- Already applied migrations are skipped

### Transaction Safety
- Each migration runs in its own database transaction
- If a migration fails, it's rolled back automatically
- Migration tracking is part of the same transaction

### SQL Statement Splitting
- Multi-statement SQL files are supported
- Statements are split by semicolons
- Empty statements and comments are ignored

### Error Handling
- Detailed error messages show which statement failed
- Failed migrations don't get marked as applied
- Application startup fails if migrations fail (fail-fast approach)

## Directory Structure

```
migrations/
├── 20241110181134_initial_schema.up.sql
├── 20241110181134_initial_schema.down.sql
├── 20241112143022_add_products.up.sql
└── 20241112143022_add_products.down.sql

scripts/
├── create_migration.sh    # Create new migration files
└── migrate.sh            # Migration management commands
```

## Best Practices

### 1. Always Create Both Files
- Create both `.up.sql` and `.down.sql` files
- Ensure rollback migrations properly undo the forward migration

### 2. Test Migrations
```bash
# Test forward migration
./scripts/migrate.sh up

# Test rollback
./scripts/migrate.sh down

# Test forward again
./scripts/migrate.sh up
```

### 3. Use Transactions
- Each migration file should be atomic
- Avoid mixing DDL and DML in the same migration if possible

### 4. Naming Conventions
- Use descriptive migration names
- Use snake_case for migration names
- Be specific: `add_user_email_index` not `update_users`

### 5. Schema Changes
- Always add new columns as nullable first
- Use separate migrations for data transformations
- Drop columns in separate migrations after data cleanup

## Common Migration Patterns

### Adding a Table
```sql
-- up.sql
CREATE TABLE new_table (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- down.sql
DROP TABLE IF EXISTS new_table;
```

### Adding a Column
```sql
-- up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- down.sql
ALTER TABLE users DROP COLUMN phone;
```

### Adding an Index
```sql
-- up.sql
CREATE INDEX idx_users_email ON users(email);

-- down.sql
DROP INDEX IF EXISTS idx_users_email;
```

### Data Migration
```sql
-- up.sql
UPDATE users SET status = 'active' WHERE status IS NULL;

-- down.sql
UPDATE users SET status = NULL WHERE status = 'active';
```

## Troubleshooting

### Migration Fails
1. Check the error message for the specific SQL statement
2. Verify your SQL syntax for PostgreSQL
3. Ensure all referenced tables/columns exist
4. Check for constraint violations

### Migration Stuck
1. Check the `migration_records` table
2. Manually remove the record if needed (be careful!)
3. Fix the migration file and retry

### Rollback Issues
1. Ensure the `.down.sql` file exists
2. Verify the rollback SQL is correct
3. Check for data dependencies that prevent rollback

## Integration with GORM

The migration system works alongside GORM:

1. **SQL migrations run first** - Handle complex schema changes
2. **GORM AutoMigrate runs second** - Sync Go models with database
3. **This ensures compatibility** between custom SQL and GORM expectations

## Security Notes

- Migration files are executed with full database privileges
- Review all migration files before applying to production
- Test migrations on a copy of production data first
- Keep migration files in version control
