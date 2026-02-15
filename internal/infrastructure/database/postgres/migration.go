package postgres

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	ID        int    `db:"id"`
	Migration string `db:"migration"`
	Batch     int    `db:"batch"`
}

func MigrateUp(db *sqlx.DB, migrationPath string) error {
	// Pastikan tabel migrations ada (Postgres syntax)
	createTable := `
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		migration VARCHAR(255) UNIQUE NOT NULL,
		batch INT
	);`
	if _, err := db.Exec(createTable); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %v", err)
	}

	// Ambil semua migration yang sudah dijalankan
	var executedMigrations []Migration
	if err := db.Select(&executedMigrations, "SELECT id, migration, batch FROM migrations"); err != nil {
		return fmt.Errorf("failed to fetch executed migrations: %v", err)
	}

	executed := make(map[string]bool)
	for _, m := range executedMigrations {
		executed[m.Migration] = true
	}

	// Ambil batch terakhir
	var lastBatch int
	if err := db.Get(&lastBatch, "SELECT COALESCE(MAX(batch), 0) FROM migrations"); err != nil {
		return fmt.Errorf("failed to get last batch: %v", err)
	}
	newBatch := lastBatch + 1

	// Baca file migrasi
	files, err := os.ReadDir(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %v", err)
	}

	// Kumpulkan dan sort migration files
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationName := strings.TrimSuffix(file.Name(), ".up.sql")
			if !executed[migrationName] {
				migrationFiles = append(migrationFiles, file.Name())
			}
		}
	}

	// Sort migration files untuk memastikan urutan yang benar
	sort.Strings(migrationFiles)

	if len(migrationFiles) == 0 {
		fmt.Println("No new migrations to run")
		return nil
	}

	// Jalankan migrations
	for _, fileName := range migrationFiles {
		migrationName := strings.TrimSuffix(fileName, ".up.sql")

		fmt.Printf("Running migration: %s\n", fileName)

		content, err := os.ReadFile(filepath.Join(migrationPath, fileName))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", fileName, err)
		}

		tx, err := db.Beginx()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}

		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute SQL statement: %v", err)
			}
		}

		if _, err := tx.Exec("INSERT INTO migrations (migration, batch) VALUES ($1, $2)", migrationName, newBatch); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration: %v", err)
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to commit migration: %v", err)
		}
	}

	return nil
}

func MigrateDown(db *sqlx.DB, migrationPath string) error {
	// Pastikan tabel migrations ada
	createTable := `
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		migration VARCHAR(255) UNIQUE NOT NULL,
		batch INT
	);`
	if _, err := db.Exec(createTable); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %v", err)
	}

	// Ambil batch terakhir
	var lastBatch int
	if err := db.Get(&lastBatch, "SELECT COALESCE(MAX(batch), 0) FROM migrations"); err != nil {
		return fmt.Errorf("failed to get last batch: %v", err)
	}

	if lastBatch == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// Ambil migrations dari batch terakhir
	var migrations []Migration
	if err := db.Select(&migrations, "SELECT id, migration, batch FROM migrations WHERE batch = $1 ORDER BY id DESC", lastBatch); err != nil {
		return fmt.Errorf("failed to fetch migrations for rollback: %v", err)
	}

	// Rollback migrations
	for _, m := range migrations {
		downFile := filepath.Join(migrationPath, m.Migration+".down.sql")

		// Check if down file exists
		if _, err := os.Stat(downFile); os.IsNotExist(err) {
			return fmt.Errorf("rollback file not found: %s", downFile)
		}

		content, err := os.ReadFile(downFile)
		if err != nil {
			return fmt.Errorf("failed to read rollback file %s: %v", downFile, err)
		}

		fmt.Printf("Rolling back: %s\n", m.Migration)

		tx, err := db.Beginx()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %v", err)
		}

		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute rollback SQL: %v", err)
			}
		}

		if _, err := tx.Exec("DELETE FROM migrations WHERE migration = $1", m.Migration); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record: %v", err)
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to commit rollback: %v", err)
		}
	}

	return nil
}

func DropAllTables(db *sqlx.DB) error {
	// This file targets Postgres by default. Attempt Postgres drop.
	return dropAllTablesPostgreSQL(db)
}

func dropAllTablesPostgreSQL(db *sqlx.DB) error {
	// Get all table names
	var tables []string
	if err := db.Select(&tables, `
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
	`); err != nil {
		return fmt.Errorf("failed to get table names: %v", err)
	}

	// Drop tables with CASCADE to handle foreign key constraints
	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, table)); err != nil {
			return fmt.Errorf("failed to drop table %s: %v", table, err)
		}
		fmt.Printf("Dropped table: %s\n", table)
	}

	return nil
}

func dropAllTablesMySQL(db *sqlx.DB) error {
	// Disable foreign key checks
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %v", err)
	}

	// Get all table names
	var tables []string
	if err := db.Select(&tables, "SHOW TABLES"); err != nil {
		return fmt.Errorf("failed to get table names: %v", err)
	}

	// Drop all tables
	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
			return fmt.Errorf("failed to drop table %s: %v", table, err)
		}
		fmt.Printf("Dropped table: %s\n", table)
	}

	// Re-enable foreign key checks
	if _, err := db.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		return fmt.Errorf("failed to re-enable foreign key checks: %v", err)
	}

	return nil
}

func dropAllTablesSQLite(db *sqlx.DB) error {
	// Get all table names
	var tables []string
	if err := db.Select(&tables, `
		SELECT name
		FROM sqlite_master
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
	`); err != nil {
		return fmt.Errorf("failed to get table names: %v", err)
	}

	// Drop all tables
	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
			return fmt.Errorf("failed to drop table %s: %v", table, err)
		}
		fmt.Printf("Dropped table: %s\n", table)
	}

	return nil
}
