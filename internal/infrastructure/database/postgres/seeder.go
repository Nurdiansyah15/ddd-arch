package postgres

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

// SeedUp executes all seeder SQL files (*.sql) in the provided directory.
func SeedUp(db *sqlx.DB, seederPath string) error {
	files, err := os.ReadDir(seederPath)
	if err != nil {
		return fmt.Errorf("failed to read seeder directory: %v", err)
	}

	var seedFiles []string
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".sql") && !strings.HasSuffix(name, ".down.sql") {
			seedFiles = append(seedFiles, name)
		}
	}

	if len(seedFiles) == 0 {
		fmt.Println("No seed files found")
		return nil
	}

	sort.Strings(seedFiles)

	for _, fileName := range seedFiles {
		fmt.Printf("Running seeder: %s\n", fileName)
		content, err := os.ReadFile(filepath.Join(seederPath, fileName))
		if err != nil {
			return fmt.Errorf("failed to read seeder file %s: %v", fileName, err)
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
				return fmt.Errorf("failed to execute seeder SQL: %v", err)
			}
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to commit seeder %s: %v", fileName, err)
		}
	}

	return nil
}

// SeedDown executes all down-seeder files (*.down.sql) in reverse order.
func SeedDown(db *sqlx.DB, seederPath string) error {
	files, err := os.ReadDir(seederPath)
	if err != nil {
		return fmt.Errorf("failed to read seeder directory: %v", err)
	}

	var downFiles []string
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".down.sql") {
			downFiles = append(downFiles, name)
		}
	}

	if len(downFiles) == 0 {
		fmt.Println("No down seeder files found")
		return nil
	}

	sort.Sort(sort.Reverse(sort.StringSlice(downFiles)))

	for _, fileName := range downFiles {
		fmt.Printf("Running seeder rollback: %s\n", fileName)
		content, err := os.ReadFile(filepath.Join(seederPath, fileName))
		if err != nil {
			return fmt.Errorf("failed to read down seeder file %s: %v", fileName, err)
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
				return fmt.Errorf("failed to execute down seeder SQL: %v", err)
			}
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to commit down seeder %s: %v", fileName, err)
		}
	}

	return nil
}
