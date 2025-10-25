package database

import (
	"fmt"
	"job-queuer/app/util/env"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() error {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.GetDBUser(),
		env.GetDBPassword(),
		env.GetDBHost(),
		env.GetDBPort(),
		env.GetDBName(),
	)

	migrationsPath := "file://./app/migrations"

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run all up migrations
	err2 := m.Up()
	if err2 != nil && err != migrate.ErrNoChange {
		log.Printf("Migration failed: %v. Attempting rollback...", err)
		if rbErr := RollbackLastMigration(); rbErr != nil {
			return fmt.Errorf("migration failed: %v; rollback also failed: %w", err, rbErr)
		}
		return fmt.Errorf("migration failed and rollback executed: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

func RollbackLastMigration() error {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.GetDBUser(),
		env.GetDBPassword(),
		env.GetDBHost(),
		env.GetDBPort(),
		env.GetDBName(),
	)

	migrationsPath := "file://./app/migrations"

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Rollback one migration
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("rollback failed: %w", err)
	}

	log.Println("Rolled back last migration")
	return nil
}
