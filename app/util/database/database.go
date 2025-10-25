package database

import (
	"database/sql"
	"fmt"
	"job-queuer/app/util/env"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	host := env.GetDBHost()
	port := env.GetDBPort()
	user := env.GetDBUser()
	password := env.GetDBPassword()
	dbname := env.GetDBName()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Run migrations using database/sql
			sqlDB, err2 := sql.Open("postgres", dsn)
			if err2 != nil {
				log.Fatalf("Failed to open database for migration: %v", err2)
			}
			err2 = RunMigrations(sqlDB, "./app/migrations")
			if err2 != nil {
				log.Fatalf("Migration failed: %v", err2)
			}
			sqlDB.Close()
			return db, nil
		}
		fmt.Printf("Failed to connect to DB (attempt %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
