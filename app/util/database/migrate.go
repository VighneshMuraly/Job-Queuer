package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s/%s", migrationsDir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}
		log.Printf("Running migration: %s", file.Name())
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}
	}
	return nil
}
