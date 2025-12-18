package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// MigrationRecord tracks which migrations have been applied
type MigrationRecord struct {
	ID        uint   `gorm:"primaryKey"`
	Filename  string `gorm:"uniqueIndex;not null"`
	AppliedAt string `gorm:"not null"`
}

// Migrator handles SQL file migrations
type Migrator struct {
	db            *gorm.DB
	migrationsDir string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// RunMigrations scans and runs pending SQL migrations
func (m *Migrator) RunMigrations() error {
	log.Println("🔄 Checking for SQL migrations...")

	// Ensure migrations table exists
	if err := m.ensureMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get list of migration files
	migrationFiles, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %v", err)
	}

	if len(migrationFiles) == 0 {
		log.Println("✅ No migration files found")
		return nil
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// Find pending migrations
	pendingMigrations := m.findPendingMigrations(migrationFiles, appliedMigrations)

	if len(pendingMigrations) == 0 {
		log.Println("✅ All migrations are up to date")
		return nil
	}

	log.Printf("📋 Found %d pending migrations", len(pendingMigrations))

	// Run pending migrations
	for _, migration := range pendingMigrations {
		if err := m.runMigration(migration); err != nil {
			return fmt.Errorf("failed to run migration %s: %v", migration, err)
		}
	}

	log.Printf("✅ Successfully applied %d migrations", len(pendingMigrations))
	return nil
}

// ensureMigrationsTable creates the migrations tracking table
func (m *Migrator) ensureMigrationsTable() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

// getMigrationFiles scans the migrations directory for .up.sql files
func (m *Migrator) getMigrationFiles() ([]string, error) {
	pattern := filepath.Join(m.migrationsDir, "*.up.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// Extract just the filenames and sort them
	var migrations []string
	for _, file := range files {
		filename := filepath.Base(file)
		migrations = append(migrations, filename)
	}

	// Sort by timestamp (filename starts with timestamp)
	sort.Strings(migrations)

	return migrations, nil
}

// getAppliedMigrations gets list of already applied migrations
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	var records []MigrationRecord
	if err := m.db.Find(&records).Error; err != nil {
		return nil, err
	}

	applied := make(map[string]bool)
	for _, record := range records {
		applied[record.Filename] = true
	}

	return applied, nil
}

// findPendingMigrations finds migrations that haven't been applied yet
func (m *Migrator) findPendingMigrations(allMigrations []string, applied map[string]bool) []string {
	var pending []string
	for _, migration := range allMigrations {
		if !applied[migration] {
			pending = append(pending, migration)
		}
	}
	return pending
}

// runMigration executes a single migration file
func (m *Migrator) runMigration(filename string) error {
	log.Printf("⚡ Running migration: %s", filename)

	// Read the migration file
	filePath := filepath.Join(m.migrationsDir, filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	sql := string(content)

	// Skip empty files or files with only comments
	if m.isEmptyMigration(sql) {
		log.Printf("⏭️  Skipping empty migration: %s", filename)
		return m.recordMigration(filename)
	}

	// Execute the SQL in a transaction
	tx := m.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	// Split SQL by semicolons and execute each statement
	statements := m.splitSQL(sql)
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		if err := tx.Exec(statement).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute statement %d in %s: %v\nStatement: %s", i+1, filename, err, statement)
		}
	}

	// Record the migration as applied
	record := MigrationRecord{
		Filename:  filename,
		AppliedAt: fmt.Sprintf("%d", m.getCurrentTimestamp()),
	}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	log.Printf("✅ Migration completed: %s", filename)
	return nil
}

// recordMigration records a migration as applied without executing SQL
func (m *Migrator) recordMigration(filename string) error {
	record := MigrationRecord{
		Filename:  filename,
		AppliedAt: fmt.Sprintf("%d", m.getCurrentTimestamp()),
	}

	return m.db.Create(&record).Error
}

// isEmptyMigration checks if a migration file is effectively empty
func (m *Migrator) isEmptyMigration(sql string) bool {
	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line != "" && !strings.HasPrefix(line, "--") {
			return false
		}
	}
	return true
}

// splitSQL splits SQL content into individual statements
func (m *Migrator) splitSQL(sql string) []string {
	// Simple split by semicolon - this might need to be more sophisticated
	// for complex SQL with semicolons in strings, etc.
	statements := strings.Split(sql, ";")
	
	var result []string
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}
	
	return result
}

// getCurrentTimestamp returns current Unix timestamp
func (m *Migrator) getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// RollbackMigration rolls back the last applied migration
func (m *Migrator) RollbackMigration() error {
	log.Println("🔄 Rolling back last migration...")

	// Get the last applied migration
	var lastRecord MigrationRecord
	if err := m.db.Order("applied_at DESC").First(&lastRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("ℹ️  No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to get last migration: %v", err)
	}

	// Find the corresponding .down.sql file
	downFilename := strings.Replace(lastRecord.Filename, ".up.sql", ".down.sql", 1)
	downFilePath := filepath.Join(m.migrationsDir, downFilename)

	// Check if down migration exists
	content, err := ioutil.ReadFile(downFilePath)
	if err != nil {
		return fmt.Errorf("rollback file not found: %s", downFilename)
	}

	sql := string(content)

	// Execute the rollback SQL
	tx := m.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin rollback transaction: %v", tx.Error)
	}

	statements := m.splitSQL(sql)
	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		if err := tx.Exec(statement).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute rollback statement %d: %v\nStatement: %s", i+1, err, statement)
		}
	}

	// Remove the migration record
	if err := tx.Delete(&lastRecord).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback: %v", err)
	}

	log.Printf("✅ Rolled back migration: %s", lastRecord.Filename)
	return nil
}
